package server

import (
	"github.com/labstack/echo"
	"github.com/kawaapp/kawaqing/model"
	"github.com/kawaapp/kawaqing/router/mwx/session"
	"github.com/kawaapp/kawaqing/store"
	"github.com/kawaapp/kawaqing/events"

	"net/http"
	"strconv"
	"fmt"
)

// likes
func GetDiscussionLikeList(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "id not found")
	}
	page, size := getPageSize(c)
	q := model.QueryParams{
		"target_ty": model.LikeDiscussion,
		"target_id": strconv.Itoa(id),
	}
	likes, err := store.GetLikeList(c, q, page-1, size)
	if err != nil {
		return fmt.Errorf("GetDiscussionLikeList, %v", err)
	}
	users, err := getLikeUserList(c, likes)
	if err != nil {
		return fmt.Errorf("GetDiscussionLikeList, %v", err)
	}
	return jsonResp(c, 0, users)
}

func GetPostLikeList(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "id not found")
	}
	page, size := getPageSize(c)
	q := model.QueryParams{
		"target_ty": model.LikePost,
		"target_id": strconv.Itoa(id),
	}
	likes, err := store.GetLikeList(c, q, page-1, size)
	if err != nil {
		return fmt.Errorf("GetPostLikeList, %v", err)
	}
	users, err := getLikeUserList(c, likes)
	if err != nil {
		return fmt.Errorf("GetPostLikeList, %v", err)
	}
	return jsonResp(c, 0, users)
}

func CreateLike(c echo.Context) error {
	in := struct {
		T string `json:"type"`
		Id int64 `json:"id"`
	}{}
	if err := c.Bind(&in); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// if user silenced
	user := session.User(c)
	if user.Silenced() {
		return c.String(http.StatusForbidden, "!user silenced")
	}

	err, firstTime := store.CreateLike(c, in.T, in.Id, user.ID)
	if err != nil {
		return fmt.Errorf("CreateLike, %v", err)
	}

	// publish favor message
	like, err := store.GetLike(c, in.T, in.Id, user.ID)
	if err != nil {
		fmt.Errorf("CreateLike, %v", err)
	}
	if firstTime {
		events.Dispatch(eLikeCreated, c, like)
	} else {
		events.Dispatch(eLikeUpdated, c, like)
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
	q := model.QueryParams{
		"user_id": strconv.Itoa(int(uid)),
	}
	favors, err := store.FromContext(c).GetLikeList(q, page, size)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if len(favors) > 0 {
		// attachActorToFavor(c, favors)
	}
	return c.JSON(200, favors)
}

// delete
func DeletePostLike(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "id not found")
	}
	user := session.User(c)
	if user.Silenced() {
		return c.String(http.StatusForbidden, "!user silenced")
	}

	db := store.FromContext(c)
	like, err := db.GetLike(model.LikePost,  int64(id), user.ID)
	if err != nil {
		return fmt.Errorf("DeletePostLike, get, %v", err)
	}

	// return if deleted
	if like.Status == 0 {
		return c.NoContent(200)
	}

	// update favor status
	if err := db.DeleteLike(model.LikePost, int64(id), user.ID); err != nil {
		return fmt.Errorf("DeletePostLike, del, %v", err)
	}
	events.Dispatch(eLikeUpdated, c, like)
	return c.NoContent(200)
}

func DeleteDiscussionLike(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "id not found")
	}
	user := session.User(c)
	if user.Silenced() {
		return c.String(http.StatusForbidden, "!user silenced")
	}

	db := store.FromContext(c)
	like, err := db.GetLike(model.LikeDiscussion, int64(id), user.ID)
	if err != nil {
		return fmt.Errorf("DeleteDiscussionLike, get, %v", err)
	}

	// return if deleted
	if like.Status == 0 {
		return c.NoContent(200)
	}

	// update favor status
	if err := db.DeleteLike(model.LikePost, int64(id), user.ID); err != nil {
		return fmt.Errorf("DeleteDiscussionLike, del, %v", err)
	}
	events.Dispatch(eLikeUpdated, c, like)
	return c.NoContent(200)
}

func getLikeUserList(c echo.Context, likes []*model.Like) ([]*model.User, error) {
	ids := make([]int64, len(likes))
	for i, v := range likes {
		ids[i] = v.UserID
	}
	return store.FromContext(c).GetUserIdList(ids)
}