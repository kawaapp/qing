package server

import (
	"github.com/labstack/echo"
	"github.com/kawaapp/kawaqing/store"
	"fmt"
)

// admin - api

func SearchUser(c echo.Context) error {
	var (
		page, size = getPageSize(c)
		params = getQueryParams(c)
	)
	db := store.FromContext(c)
	users, err := db.GetUserList(params, page-1, size)
	if err != nil {
		return fmt.Errorf("SearchUser, %v", err)
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
	db := store.FromContext(c)
	discussions, err := db.GetDiscussionList(params, page-1, size)
	if err != nil {
		return fmt.Errorf("SearchDiscussions, %v", err)
	}

	p := makePayload(0, discussions)
	if page == 0 {
		num, _ := db.GetDiscussionCount(params)
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
	db := store.FromContext(c)
	posts, err := db.GetPostList(params, page-1, size)
	if err != nil {
		return fmt.Errorf("SearchPosts, %v", err)
	}
	p := makePayload(0, posts)
	if page == 0 {
		num, _ := db.GetPostCount(params)
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
	db := store.FromContext(c)
	reports, err := db.GetReportList(params, page-1, size)
	if err != nil {
		return fmt.Errorf("SearchReport, %v", err)
	}
	p := makePayload(0, reports)
	if page == 0 {
		num, _ := db.GetReportCount(params)
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