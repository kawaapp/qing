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
	"log"
)

// discussion
type discussion struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	CateID  int64 `json:"cate_id"`
}

type outDiscussion struct {
	*model.Discussion
	Liked bool `json:"liked"`
	Favorited bool `json:"favorited"`
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
		sort = c.QueryParam("q")
		//sort   = c.FormValue("sort")
		page, size = getPageSize(c)
	)
	db := store.FromContext(c)
	q := model.QueryParams{
		"sort": sort,
	}
	if cid := c.QueryParam("cate_id"); len(cid) > 0 {
		q["cate_id"] = cid
	}

	discussions, err := db.GetDiscussionList(q, page-1, size)
	if err != nil {
		return fmt.Errorf("GetDiscussionList, %v", err)
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
	if err == sql.ErrNoRows {
		return c.NoContent(404)
	}
	if err != nil {
		return fmt.Errorf("GetDiscussion, %v", err)
	}

	// increase counter
	go increaseViewCount(c, discussion)

	// get user like state
	var (
		liked = false
		favor = false
	)

	// optional user
	if usr := session.User(c); usr != nil {
		_, err := store.GetLike(c, "dz", int64(id), usr.ID)
		if err == nil {
			liked = true
		}
		_, err = store.FromContext(c).GetFavoriteUser(usr.ID, int64(id))
		if err == nil {
			favor = true
		}
	}

	// with extra state
	out := outDiscussion{
		Discussion: discussion, Liked: liked, Favorited: favor,
	}

	p := makePayload(0, out)

	// attach users
	if includes(c, "user") {
		attachUserToDiscussion(c, []*model.Discussion{ discussion }, p)
	}
	return c.JSON(200, p)
}

func increaseViewCount(c echo.Context, d *model.Discussion) {
	d.ViewCount += 1
	err := store.UpdateDiscussion(c, d)
	if err != nil {
		log.Printf("increaseViewCount, %v", err)
	}
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
		CategoryID: in.CateID,
	}
	if err := store.CreateDiscussion(c, discussion); err != nil {
		return fmt.Errorf("CreateDiscussion, %v", err)
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
		return fmt.Errorf("AdminCreateDiscussion, %v", err)
	}
	return c.JSON(200, discussion)
}

func UpdateDiscussion(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	d, err := store.FromContext(c).GetDiscussion(int64(id))
	if err != nil {
		return fmt.Errorf("UpdateDiscussion, %v", err)
	}

	// ensure author is valid
	user := session.User(c)
	if user.ID != d.AuthorID {
		return c.NoContent(http.StatusForbidden)
	}

	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	in := (map[string]interface{})(m)

	// partial update
	if value, ok := in["title"]; ok {
		d.Title = value.(string)
	}
	if value, ok := in["content"]; ok {
		d.Content = value.(string)
	}
	if err := store.UpdateDiscussion(c, d); err != nil {
		return fmt.Errorf("UpdateDiscussion, %v", err)
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
		return fmt.Errorf("DeleteDiscussion, %v", err)
	}
	// ensure author is valid
	if user := session.User(c); user.Login != "Admin" && user.ID != d.AuthorID {
		return c.NoContent(http.StatusForbidden)
	}

	if err := store.DeleteDiscussion(c, int64(id)); err != nil {
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
		return fmt.Errorf("GetDiscussionByUser, %v", err)
	}

	page, size := getPageSize(c)
	q := model.QueryParams {
		"author_id": strconv.Itoa(int(user.ID)),
	}
	discussions, err := store.FromContext(c).GetDiscussionList(q, page-1, size)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	p := makePayload(0, discussions)

	// attach users
	if includes(c, "user") {
		p.Entities["users"] = map[int64]*model.User {
			user.ID: user,
		}
	}
	return c.JSON(200, p)
}

func attachUserToDiscussion(c echo.Context, posts []*model.Discussion, p payload) {
	ids := make([]int64, len(posts) * 2)
	kv := make(map[int64]*model.User)

	// attach user info
	for i, v := range posts {
		ids[i<<1+0] = v.AuthorID
		ids[i<<1+1] = v.LastReplyUid
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
	p, ok := v.(*model.Post)
	if !ok {
		return typeError("Post")
	}
	d, err := store.FromContext(c).GetDiscussion(p.DiscussionID)
	if err != nil {
		return err
	}
	{
		d.CommentCount = getCount(d.CommentCount)
		d.LastReplyAt = p.CreatedAt
		d.LastReplyUid = p.AuthorID
	}
	err = store.FromContext(c).UpdateDiscussion(d)
	if err != nil {
		return fmt.Errorf("dzOnCommentChanged, %v", err)
	}
	return nil
}

func dzOnCommentCreated(c echo.Context, v interface{}) error  {
	return dzOnCommentChanged(c, v, func(num int) int {
		return num + 1
	})
}

func dzOnCommentDeleted(c echo.Context, v interface{}) error  {
	return dzOnCommentChanged(c, v, func(num int) int {
		return Max(num -1, 0)
	})
}

func dzOnLikeCreated(c echo.Context, v interface{}) error {
	return dzOnLikeChanged(c, v, func(base int) int {
		return base + 1
	})
}

func dzOnLikeDeleted(c echo.Context, v interface{}) error {
	return dzOnLikeChanged(c, v, func(base int) int {
		return Max(base-1, 0)
	})
}

func dzOnLikeChanged(c echo.Context, v interface{}, getCount func(base int) int) error {
	like, ok := v.(*model.Like)
	if !ok {
		return typeError("Like")
	}
	if like.TargetTy != model.LikeDiscussion {
		return nil  // skip if not dz
	}
	d, err := store.GetDiscussion(c, like.TargetID)
	if err != nil {
		return err
	}
	d.LikeCount = getCount(d.LikeCount)
	err = store.UpdateDiscussion(c, d)
	if err != nil {
		return fmt.Errorf("dzOnLikeChanged, %v", err)
	}
	return nil
}

func init() {
	events.Subscribe(ePostCreated, dzOnCommentCreated)
	events.Subscribe(ePostDeleted, dzOnCommentDeleted)
	events.Subscribe(eLikeCreated, dzOnLikeCreated)
	events.Subscribe(eLikeDeleted, dzOnLikeDeleted)
}
