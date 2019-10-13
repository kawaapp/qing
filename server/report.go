package server

import (
	"github.com/labstack/echo"
	"github.com/kawaapp/kawaqing/model"
	"github.com/kawaapp/kawaqing/router/mwx/session"
	"github.com/kawaapp/kawaqing/store"

	"net/http"
	"strconv"
	"database/sql"
)

func CreateReport(c echo.Context) error  {
	report := new(model.Report)
	if err := c.Bind(report); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	report.Counter = 1
	report.UserId = session.User(c).ID

	err := store.FromContext(c).CreateReport(report)
	if err != nil {
		return err
	}
	return c.JSON(200, report)
}

func GetReport(c echo.Context) error  {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	rpt, err := store.FromContext(c).GetReport(int64(id))
	if err == sql.ErrNoRows {
		return c.NoContent(http.StatusNotFound)
	}
	if err != nil {
		return err
	}
	return c.JSON(200, rpt)
}

func SetReportStatus(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "id not found")
	}
	st, err := strconv.Atoi(c.QueryParam("st"))
	if err != nil {
		return c.String(http.StatusBadRequest, "status not found")
	}
	rpt, err := store.FromContext(c).GetReport(int64(id))
	if err == sql.ErrNoRows {
		return c.NoContent(404)
	}
	if err != nil {
		return err
	}
	rpt.Status = int64(st)
	err = store.FromContext(c).UpdateReport(rpt)
	if err != nil {
		return err
	}
	return c.JSON(200, rpt)
}

func attachUserToReport(c echo.Context, reports []*model.Report, p payload) {
	ids := make([]int64, len(reports))
	kv := make(map[int64]*model.User)

	// attach user info
	for i, v := range reports {
		ids[i] = v.UserId
	}

	// 返回的user-id可能会比post少，sql查询是去重的结果
	authors, _ := store.FromContext(c).GetUserIdList(ids)
	for _, v := range authors {
		kv[v.ID] = v
	}
	p.Entities["users"] = kv
}

func attachPostToReport(c echo.Context, reports []*model.Report, p payload) {
	ids := make([]int64, len(reports))
	kv := make(map[int64]*model.Post)

	// attach post info
	for i, v := range reports {
		ids[i] = v.EntityId
	}

	posts, _ := store.FromContext(c).GetPostListByIds(ids)
	for _, v := range posts {
		kv[v.ID] = v
	}
	p.Entities["posts"] = kv
}