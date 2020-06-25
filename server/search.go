package server

import (
	"github.com/labstack/echo"
	"github.com/kawaapp/kawaqing/store"
)

// admin - api

func SearchUser(c echo.Context) error {
	var (
		page, size = getPageSize(c)
		params = getQueryParams(c)
	)
	db := store.FromContext(c)
	users, err := db.GetUserList(params, page, size)
	if err != nil {
		return err
	}

	p := makePayload(0, users)
	if page == 0 {
		num, _ := db.GetUserCount(params)
		p.Total = num
	}
	return c.JSON(200, p)
}

func SearchDiscussions(c echo.Context) error {
	var (
		page, size = getPageSize(c)
		params = getQueryParams(c)
	)
	discussions, err := store.SearchDiscussion(c, params, page, size)
	if err != nil {
		return err
	}

	p := makePayload(0, discussions)
	if page == 0 {
		num, _ := store.SearchDiscussionCount(c, params)
		p.Total = num
	}

	if includes(c, "user") {
		attachUserToDiscussion(c, discussions, p)
	}
	return c.JSON(200, p)
}

func SearchPosts(c echo.Context) error {
	var (
		page, size = getPageSize(c)
		params = getQueryParams(c)
	)
	posts, err := store.SearchPost(c, params, page, size)
	if err != nil {
		return err
	}
	p := makePayload(0, posts)
	if page == 0 {
		num, _ := store.SearchPostCount(c, params)
		p.Total = num
	}
	if includes(c, "user") {
		attachUserToPost(c, posts, p)
	}

	if includes(c, "like") {
		attachLikeToPost(c, posts, p)
	}
	return c.JSON(200, p)
}

func SearchReport(c echo.Context) error {
	var (
		page, size = getPageSize(c)
		params = getQueryParams(c)
	)
	reports, err := store.SearchReport(c, params, page, size)
	if err != nil {
		return err
	}
	p := makePayload(0, reports)
	if page == 0 {
		num, _ := store.SearchReportCount(c, params)
		p.Total = num
	}
	if includes(c, "user") {
		attachUserToReport(c, reports, p)
	}
	if includes(c, "post") {
		attachPostToReport(c, reports, p)
	}
	return c.JSON(200, p)
}

func SearchSignUser(c echo.Context) error {
	var (
		page = GetQueryIntValue(c.QueryParam("page"), 1)
		size = GetQueryIntValue(c.QueryParam("size"), 30)
	)
	pager, err := store.FromContext(c).SearchSignUser(page, size)
	if err != nil {
		return err
	}
	return c.JSON(200, pager)
}