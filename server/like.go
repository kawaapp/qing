package server

import (
	"github.com/labstack/echo"
	"github.com/kawaapp/kawaqing/model"
	"github.com/kawaapp/kawaqing/router/mwx/session"
	"github.com/kawaapp/kawaqing/store"
	"github.com/kawaapp/kawaqing/events"

	"net/http"
	"strconv"
	"database/sql"
)

// likes
func GetLikeList(c echo.Context) error {
	cid, err := strconv.Atoi(c.Param("cid"))
	if err != nil {
		return c.String(http.StatusBadRequest, "cid not found")
	}
	return getFavorList(c, int64(cid))
}

func GetLikeCount(c echo.Context) error {
	cid, err := strconv.Atoi(c.Param("cid"))
	if err != nil {
		return c.String(http.StatusBadRequest, "cid not found")
	}
	return getFavorCount(c, int64(cid))
}

func CreateLike(c echo.Context) error {
	in := struct {
		Cid int64 `json:"id"`
	}{}
	if err := c.Bind(&in); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// if user silenced
	user := session.User(c)
	if user.Silenced() {
		return c.String(http.StatusForbidden, "!user silenced")
	}

	_, err := store.GetPost(c, int64(in.Cid))
	if err != nil {
		return c.String(http.StatusBadRequest, "comment not exist")
	}

	var firstTime bool
	f, err := store.FromContext(c).GetLike(in.Cid, user.ID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	// return if already favored
	if err == nil && f.Status == 1 {
		return c.NoContent(200)
	}

	// update if exist
	if err == nil && f.Status == 0 {
		f.Status = 1
		err = store.FromContext(c).UpdateLike(f)
	}

	// create if not exist
	if err == sql.ErrNoRows {
		f = &model.Like{
			Status:     1,
			AuthorID:   session.User(c).ID,
			PostId: in.Cid,
		}
		err = store.CreateFavor(c, f)
		firstTime = true
	}

	if err != nil {
		return err
	}

	// publish favor message
	if firstTime {
		events.Dispatch(eLikeCreated, c, f)
	} else {
		events.Dispatch(eLikeUpdated, c, f)
	}

	return c.NoContent(200)
}

func GetLikeByUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	uid := int64(id)
	if id == 0 {
		uid = session.User(c).ID
	}
	page, size := getPageSize(c)
	favors, err := store.FromContext(c).GetLikeListUser(uid, page, size)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if len(favors) > 0 {
		// attachActorToFavor(c, favors)
	}
	return c.JSON(200, favors)
}

// WARNING: it's the post's author, not the one give the favor.
//func attachActorToFavor(c echo.Context, favors []*model.Like) {
//	Ids := make([]int64, len(favors))
//	kv := make(map[int64]*model.User)
//
//	for i, v := range favors {
//		Ids[i] = v.ActorID
//	}
//	from, _ := store.FromContext(c).GetUserIdList(Ids)
//	if from != nil {
//		for _, v := range from {
//			kv[v.ID] = v
//		}
//	}
//	for _, v := range favors {
//		v.Actor = kv[v.ActorID]
//	}
//}

// implementation
func DeleteLike(c echo.Context) error {
	cid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "id not found")
	}
	user := session.User(c)
	if user.Silenced() {
		return c.String(http.StatusForbidden, "!user silenced")
	}

	f, err := store.FromContext(c).GetLike(int64(cid), user.ID)
	if err == sql.ErrNoRows {
		return c.NoContent(404)
	} else if err != nil {
		return err
	}

	// return if deleted
	if f.Status == 0 {
		return c.NoContent(200)
	}

	// update favor status
	f.Status = 0
	err = store.FromContext(c).UpdateLike(f)
	if err != nil {
		return err
	}
	events.Dispatch(eLikeUpdated, c, f)
	return c.NoContent(200)
}

func getFavorList(c echo.Context, pid int64) error {
	page, size := getPageSize(c)
	favors, err := store.GetLikeList(c, pid, page, size)
	if err != nil {
		return err
	}
	return c.JSON(200, favors)
}

func getFavorCount(c echo.Context, pid int64) error {
	num, err := store.GetFavorCount(c, pid)
	if err != nil {
		return err
	}
	return c.JSON(200, struct {
		Num int `json:"num"`
	}{num})
}

