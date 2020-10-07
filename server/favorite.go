package server

import (
	"github.com/labstack/echo"
	"github.com/kawaapp/kawaqing/router/mwx/session"
	"github.com/kawaapp/kawaqing/store"
	"strconv"
	"database/sql"
	"github.com/kawaapp/kawaqing/model"
	"fmt"
)

func GetFavoriteByUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(400, "uid not found")
	}
	uid := int64(id)
	if uid == 0 {
		uid = session.User(c).ID
	}
	if uid != session.User(c).ID {
		return c.String(403, "invalid access to favorites")
	}
	page, size := getPageSize(c)

	q := model.QueryParams{
		"user_id": c.Param("id"),
	}
	favors, err := store.FromContext(c).GetFavoriteList(q, page-1, size)
	if err != nil {
		return fmt.Errorf("GetFavoriteByUser, %v", err)
	}
	dzs, err := favorsToDiscussions(c, favors)
	if err != nil {
		return fmt.Errorf("GetFavoriteByUser, %v", err)
	}

	p := makePayload(0, dzs)

	// attach users
	if includes(c, "user") {
		attachUserToDiscussion(c, dzs, p)
	}
	return c.JSON(200, p)
}

func CreateFavorite(c echo.Context) error  {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(400, "pid not found")
	}
	uid := session.User(c).ID
	did := int64(id)

	// if exist
	f, err := store.FromContext(c).GetFavoriteUser(uid, did)
	if err == nil {
		return c.JSON(200, f)
	}

	// else create new
	f.UserID = uid
	f.DiscussionID = did
	err = store.FromContext(c).CreateFavorite(f)
	if err != nil {
		return err
	}
	return c.JSON(200, f)
}

func DeleteFavorite(c echo.Context) error  {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(400, "pid not found")
	}
	uid := session.User(c).ID
	did := int64(id)

	f, err := store.FromContext(c).GetFavoriteUser(uid, did)

	// if not exist
	if err == sql.ErrNoRows {
		return c.NoContent(404)
	}

	// delete it
	err = store.FromContext(c).DeleteFavorite(f.ID)
	if err != nil {
		return err
	}
	return c.NoContent(200)
}

func favorsToDiscussions(c echo.Context, favors []*model.Favorite) ([]*model.Discussion, error) {
	ids := make([]int64, len(favors))
	for i, v := range favors {
		ids[i] = v.DiscussionID
	}
	return store.FromContext(c).GetDiscussionListByIds(ids)
}

