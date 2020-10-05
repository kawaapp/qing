package server

import (
	"github.com/labstack/echo"
	"github.com/kawaapp/kawaqing/model"
	"github.com/kawaapp/kawaqing/store"
	"github.com/kawaapp/kawaqing/events"
	"strconv"
	"github.com/kawaapp/kawaqing/router/mwx/session"
	"database/sql"
	"fmt"
)

func GetFollowerList(c echo.Context) error  {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(400, "uid not found")
	}
	uid := int64(id)
	if uid == 0 {
		uid = session.User(c).ID
	}
	page, size := getPageSize(c)
	users, err := store.FromContext(c).GetFollowerList(uid, page, size)
	if err != nil {
		return fmt.Errorf("GetFollowerList, %v", err)
	}
	return  jsonResp(c, 0, users)
}

func GetFollowingList(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(400, "uid not found")
	}
	uid := int64(id)
	if uid == 0 {
		uid = session.User(c).ID
	}
	page, size := getPageSize(c)
	users, err := store.FromContext(c).GetFollowingList(uid, page, size)
	if err != nil {
		return fmt.Errorf("GetFollowingList, %v", err)
	}
	return jsonResp(c, 0, users)
}

func GetFollowing(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("uid"))
	if err != nil {
		return c.String(400, "uid not found")
	}
	var (
		uid = int64(id)
		fid = session.User(c).ID
	)
	_, err = store.FromContext(c).GetFollow(uid, fid)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("GetFollowing, %v", err)
	}
	if err == sql.ErrNoRows {
		return c.NoContent(404)
	} else {
		return c.NoContent(200)
	}
}

func CreateFollow(c echo.Context) error {
	f := new(model.Follow)
	if err := c.Bind(f); err != nil {
		return c.String(400, err.Error())
	}
	f.FollowerId = session.User(c).ID
	err := store.FromContext(c).CreateFollow(f)
	if err != nil {
		return fmt.Errorf("CreateFollow, %v", err)
	}
	events.Dispatch(eUserFollow, c, f)
	return c.NoContent(200)
}

func DeleteFollow(c echo.Context) error  {
	id, err := strconv.Atoi(c.Param("uid"))
	if err != nil {
		return c.String(400, "uid not found")
	}
	var (
		uid = int64(id)
		fid = session.User(c).ID
	)

	f, err := store.FromContext(c).GetFollow(uid, fid)
	if err == sql.ErrNoRows {
		return c.NoContent(200)
	}
	if err != nil {
		return fmt.Errorf("DeleteFollow, %v", err)
	}
	err = store.FromContext(c).DeleteFollow(uid, fid)
	if err != nil {
		return fmt.Errorf("DeleteFollow, %v", err)
	}
	events.Dispatch(eUserUnfollow, c, f)
	return c.NoContent(200)
}
