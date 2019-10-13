package datasource

import (
	"github.com/kawaapp/kawaqing/model"
	"github.com/russross/meddler"
)

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