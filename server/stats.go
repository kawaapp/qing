package server

import (
	"github.com/labstack/echo"
	"github.com/kawaapp/kawaqing/model"
	"github.com/kawaapp/kawaqing/store"
	"time"
	"net/http"
)

// Stats overview filter by post/user/comments
func GetStats(c echo.Context) error {
	filter := c.QueryParam("q")
	if len(filter) == 0 {
		return c.String(http.StatusBadRequest, "no query params")
	}
	from, err := time.Parse("2006-01-02", c.QueryParam("from"))
	if err != nil {
		return c.String(http.StatusBadRequest, "date from:" + err.Error())
	}
	to, err   := time.Parse("2006-01-02", c.QueryParam("to"))
	if err != nil {
		return c.String(http.StatusBadRequest, "date to:" + err.Error())
	}
	to = to.Add(time.Hour*24)
	if now := time.Now(); to.After(now) {
		to = now
	}

	var (
		data []*model.DailyCount
	)
	switch filter {
	case "active-user":
		data, err = store.GetDailyActiveUser(c, from, to)
	case "new-user":
		data, err = store.GetDailyNewUser(c, from, to)
	case "new-post":
		data, err = store.GetDailyNewDiscussion(c, from, to)
	case "new-comment":
		data, err = store.GetDailyNewComment(c, from, to)
	case "new-favor":
		data, err = store.GetDailyNewFavor(c, from, to)
	}
	if err != nil {
		return err
	}
	return c.JSON(200, data)
}

func GetStatsOverView(c echo.Context) error {
	var (
		store = store.FromContext(c)
		today = time.Now()
		ov = model.StatsOverview{}
	)

	// fetch data
	ov.NewDiscussion, _ = store.GetNewDiscussion(today)
	ov.NewUser, _ = store.GetNewUser(today)
	ov.ActiveUser, _ = store.GetUserActive(today)
	ov.TotalUser, _ = store.GetTotalUser()

	// return
	return c.JSON(200, ov)
}