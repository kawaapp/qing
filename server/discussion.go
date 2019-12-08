package server

import (
	"github.com/labstack/echo"
	"github.com/kawaapp/kawaqing/model"
	"github.com/kawaapp/kawaqing/store"
	"github.com/kawaapp/kawaqing/router/mwx/session"
	"github.com/kawaapp/kawaqing/events"

	"errors"
	"fmt"
	"net/http"
	"strconv"
	"database/sql"
	"strings"
)

// discussion
type discussion struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (in *discussion) validate() error  {
	if len(in.Title) == 0 {
		return errors.New("title field required")
	}
	if len(in.Content) == 0 {
		return errors.New("content field required")
	}
	return nil
}

// user - api
func GetDiscussionList(c echo.Context) error {
	var (
		filter = c.FormValue("filter")
		//sort   = c.FormValue("sort")
		page, size = getPageSize(c)
	)

	discussions, err := store.GetDiscussionList(c, page, size, filter)
	if err != nil {
		return err
	}

	p := makePayload(0, discussions)
	if len(discussions) == size {
		p.HasMore = true
	}

	// attach users
	if includes(c, "user") {
		attachUserToDiscussion(c, discussions, p)
	}
	return c.JSON(200, p)
}

func GetDiscussion(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	discussion, err := store.GetDiscussion(c, int64(id))
	if err != nil {
		return c.NoContent(404)
	}

	p := makePayload(0, discussion)

	// attach users
	if includes(c, "user") {
		attachUserToDiscussion(c, []*model.Discussion{ discussion }, p)
	}
	return c.JSON(200, p)
}

func CreateDiscussion(c echo.Context) error {
	in := &discussion{}
	if err := c.Bind(in); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := in.validate(); err != nil {
		return c.JSON(http.StatusBadRequest, errPayload(422, err))
	}

	// ensure author is valid
	user := session.User(c)
	if user.Silenced() {
		return c.String(http.StatusForbidden, "!user silenced")
	}

	//// audit before publish
	//if ShouldAuditBeforePost(c) {
	//	post.Status = SetBitN(post.Status, 3, 1)
	//}
	//
	//// spam test
	//if checker := spamcheck.FromContext(c); checker != nil {
	//	if ok, _ := checker.Validate(post.Content); !ok {
	//		post.Status = SetBitN(post.Status, 3, 1)
	//	}
	//}

	discussion := &model.Discussion{
		AuthorID: user.ID,
		Title: in.Title,
		Content: in.Content,
	}
	if err := store.CreateDiscussion(c, discussion); err != nil {
		return err
	}
	post := &model.Post{
		DiscussionID: discussion.ID,
		AuthorID: user.ID,
		Content: in.Content,
	}
	if err := store.CreatePost(c, post); err != nil {
		return err
	}

	events.Dispatch(eDiscussionCreated, c, discussion)
	return c.JSON(200, discussion)
}

func AdminCreateDiscussion(c echo.Context) error {
	discussion := new(model.Discussion)
	if err := c.Bind(discussion); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := store.CreateDiscussion(c, discussion); err != nil {
		return err
	}
	return c.JSON(200, discussion)
}

func UpdateDiscussion(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	// ensure author is valid
	user := session.User(c)
	if user.Login != "Admin" {
		return c.NoContent(http.StatusForbidden)
	}

	d := new(model.Discussion)
	if err := c.Bind(d); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	d.ID = int64(id)
	if err := store.UpdateDiscussion(c, d); err != nil {
		return errors.New(fmt.Sprintf("Error: update post %d. %s", id, err))
	}
	return c.JSON(200, d)
}

func DeleteDiscussion(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	d, err := store.GetDiscussion(c, int64(id));
	if err == sql.ErrNoRows {
		return c.NoContent(404)
	} else if err != nil {
		return err
	}
	// ensure author is valid
	if user := session.User(c); user.Login != "Admin" && user.ID != d.AuthorID {
		return c.NoContent(http.StatusForbidden)
	}

	if err := store.DeletePost(c, int64(id)); err != nil {
		return errors.New(fmt.Sprintf("Error: delete post %d. %s", id, err))
	}
	events.Dispatch(eDiscussionDeleted, c, d)
	return c.NoContent(200)
}

func GetDiscussionByUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	var user *model.User
	if id == 0 {
		user = session.User(c)
	} else {
		user, err = store.GetUser(c, int64(id))
	}
	if err == sql.ErrNoRows {
		return c.String(404, "no user found")
	} else if err != nil {
		return err
	}

	page, size := getPageSize(c)
	discussions, err := store.FromContext(c).GetDiscussionListUser(user.ID, page, size)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	p := makePayload(0, discussions)

	// attach users
	if strings.Contains(c.QueryParam("includes"), "user") {
		p.Entities["users"] = map[int64]*model.User {
			user.ID: user,
		}
	}
	return c.JSON(200, p)
}

func attachUserToDiscussion(c echo.Context, posts []*model.Discussion, p payload) {
	ids := make([]int64, len(posts))
	kv := make(map[int64]*model.User)

	// attach user info
	for i, v := range posts {
		ids[i] = v.AuthorID
	}

	// 返回的user-id可能会比post少，sql查询是去重的结果
	authors, _ := store.FromContext(c).GetUserIdList(ids)
	for _, v := range authors {
		kv[v.ID] = v
	}
	p.Entities["users"] = kv
}


// 设置置顶/精华/推荐 form:
// top= 0, 1; 1置顶, 0 取消
// val= 0, 1; 1精华, 0 取消
// hid= 0, 1; 1隐藏, 0 取消
// aud= 0, 1; 1审核, 0 取消
func SetDiscussionStatus(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	p, err := store.GetPost(c, int64(id))
	if err != nil {
		return c.String(404, err.Error())
	}

	//if v, err := strconv.Atoi(c.FormValue("top")); (err == nil) && (v == 0 || v == 1) {
	//	p.Status = SetBitN(p.Status, 0, uint64(v))
	//}
	//if v, err := strconv.Atoi(c.FormValue("val")); (err == nil) && (v == 0 || v == 1) {
	//	p.Status = SetBitN(p.Status, 1, uint64(v))
	//}
	//if v, err := strconv.Atoi(c.FormValue("hid")); (err == nil) && (v == 0 || v == 1) {
	//	p.Status = SetBitN(p.Status, 2, uint64(v))
	//}
	//if v, err := strconv.Atoi(c.FormValue("aud")); (err == nil) && (v == 0 || v == 1) {
	//	p.Status = SetBitN(p.Status, 3, uint64(v))
	//}

	err = store.UpdatePost(c, p)
	if err != nil {
		return err
	}
	return c.JSON(200, p)
}

// events
func dzOnCommentChanged(c echo.Context, v interface{}, getCount func(num int) int) error {
	comment, ok := v.(*model.Post)
	if !ok {
		return typeError("Comment")
	}
	d, err := store.FromContext(c).GetDiscussion(comment.DiscussionID)
	if err != nil {
		return err
	}
	d.CommentCount = getCount(d.CommentCount)
	err = store.FromContext(c).UpdateDiscussion(d)
	return err
}

func postOnCommentCreated(c echo.Context, v interface{}) error  {
	return dzOnCommentChanged(c, v, func(num int) int {
		return num + 1
	})
}

func postOnCommentDeleted(c echo.Context, v interface{}) error  {
	return dzOnCommentChanged(c, v, func(num int) int {
		return Max(num -1, 0)
	})
}


func init() {
	events.Subscribe(ePostCreated, postOnCommentCreated)
	events.Subscribe(ePostDeleted, postOnCommentDeleted)
}
