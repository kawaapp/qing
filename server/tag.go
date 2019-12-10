package server

import (
	"github.com/labstack/echo"
	"github.com/kawaapp/kawaqing/store"
	"net/http"
	"strconv"
	"net/url"
	"database/sql"
	"sort"
)

// tags
func GetDiscussionsByTag(c echo.Context) error {
	tag := c.Param("tag")
	if len(tag) == 0 {
		return c.String(http.StatusBadRequest, "no tag found")
	}
	decoded, err := url.PathUnescape(tag)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	page, limit := getPageSize(c)
	discussions, err := store.GetDiscussionByTag(c, decoded, page, limit)
	if err != nil {
		return err
	}
	p := makePayload(0, discussions)
	if includes(c, "user") {
		attachUserToDiscussion(c, discussions, p)
	}
	return c.JSON(200, p)
}

func GetTagList(c echo.Context) error {
	tags, err := store.GetTagList(c)
	if err != nil {
		return err
	}
	sort.SliceStable(tags, func(i, j int) bool {
		return tags[i].Order < tags[j].Order
	})
	return c.JSON(200, tags)
}

func LinkTagPost(c echo.Context) error {
	in := struct {
		Tags []string `json:"tags"`
		PID  int      `json:"pid"`
	}{}
	if err := c.Bind(&in); err != nil {
		return err
	}
	err := store.LinkTagPost(c, int64(in.PID), in.Tags)
	if err != nil {
		return err
	}
	return c.NoContent(200)
}

func CreateTag(c echo.Context) error {
	in := struct {
		Tag string `json:"tag"`
		Summary string `json:"summary"`
	}{}
	if err := c.Bind(&in); err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}
	t, err := store.FromContext(c).CreateTag(in.Tag, in.Summary)
	if err != nil {
		return err
	}
	return c.JSON(200, t)
}

func UpdateTag(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	tag, err := store.FromContext(c).GetTagId(int64(id))
	if err == sql.ErrNoRows {
		return c.NoContent(404)
	}
	if err != nil {
		return err
	}
	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	in := (map[string]interface{})(m)

	// Override existing value, only!
	if value, ok := in["order"]; ok {
		tag.Order = int(value.(float64))
	}
	if value, ok := in["text"]; ok {
		tag.Text = value.(string)
	}
	if value, ok := in["summary"]; ok {
		tag.Summary = value.(string)
	}
	err = store.FromContext(c).UpdateTag(tag)
	if err != nil {
		return err
	}
	return c.JSON(200, tag)
}

func DeleteTag(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	err = store.FromContext(c).DeleteTag(int64(id))
	if err != nil {
		return err
	}
	return c.NoContent(200)
}