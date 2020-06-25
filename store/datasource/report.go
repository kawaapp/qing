package datasource

import (
	"github.com/kawaapp/kawaqing/model"
	"github.com/russross/meddler"
	"fmt"
)

func (db *datasource) GetReportList(params model.QueryParams, page, size int) ([]*model.Report, error) {
	reports := make([]*model.Report, 0)
	query, args := sqlReportQuery(sqlReportBase, params, page, size)
	err := meddler.QueryAll(db, &reports, query, args...)
	return reports, err
}

func (db *datasource) GetReportCount(params model.QueryParams) (int, error) {
	query, args := sqlReportQuery(sqlReportCount, params, 0, 0)
	return Count(db, query, args...)
}

func (db *datasource) CreateReport(rpt *model.Report) error {
	rpt.CreatedAt = UnixNow()
	return meddler.Insert(db, "reports", rpt)
}

func (db *datasource) UpdateReport(rpt *model.Report) error {
	rpt.UpdatedAt = UnixNow()
	return meddler.Update(db, "reports", rpt)
}

func (db *datasource) DeleteReport(id int64) error {
	_, err := db.Exec(sqlDeleteReport, id)
	return err
}

func (db *datasource) GetReport(id int64) (*model.Report, error)  {
	rpt := new(model.Report)
	err := meddler.Load(db, "reports", rpt,  id)
	return rpt, err
}

const sqlDeleteReport = `
DELETE FROM reports WHERE id=?
;`


const sqlReportBase = `
SELECT
	c.id,
	c.created_at,
	c.updated_at,
	c.entity_id,
	c.entity_ty,
	c.content,
	c.counter,
	c.status,
	c.user_id,
	c.report_ty,
	c.other,
	c.images
`

const sqlReportCount = `
SELECT COUNT(c.id)
`

func sqlReportQuery(queryBase string, params model.QueryParams, page, size int) (query string, args []interface{}) {
	query += queryBase
	query += " FROM reports c"

	// join
	query += " LEFT JOIN posts a ON a.id = c.entity_id"
	query += " LEFT JOIN users b ON b.id = c.user_id"

	// where
	where := ""
	if q, ok := params["post"]; ok {
		where += " AND a.content LIKE ?"
		args = append(args, "%" + q + "%")
	}
	if q, ok := params["user"]; ok {
		where += " AND b.nickname LIKE ?"
		args = append(args, "%" + q + "%")
	}
	if _, ok := params["status"]; ok {
		where += " AND c.status > 0"
	}

	if len(where) > 0 {
		query += " WHERE 1=1" + where
	}

	if size > 0 {
		query += fmt.Sprintf(" ORDER BY c.id DESC LIMIT %d OFFSET %d", size, page * size)
	}
	return
}