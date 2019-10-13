package datasource

import (
	"testing"
	"github.com/kawaapp/kawaqing/model"
)

func TestReportCreate(t *testing.T)  {
	s := beforeReport()
	defer s.Close()

	rpt := &model.Report{
		ReportTy: 1,
		UserId: 1,
		EntityId: 2,
	}
	if err := s.CreateReport(rpt); err != nil {
		t.Error(err)
	}

	get, err := s.GetReport(rpt.ID)
	if err != nil {
		t.Error(err)
	}
	if get.UserId != rpt.UserId {
		t.Error("reports err, get:", get," expected:", rpt)
	}


	if get.Other != rpt.Other {
		t.Error("reports err, get:", get," expected:", rpt)
	}

	if err := s.DeleteReport(rpt.ID); err != nil {
		t.Error(err)
	}

	_, err = s.GetReport(rpt.ID)
	if err == nil {
		t.Error("Report should be deleted.")
	}
}

func TestReportUpdate(t *testing.T)  {
	s := beforeReport()
	defer s.Close()

	rpt := &model.Report{
		ReportTy: 1,
		UserId: 1,
		EntityId: 2,
	}
	if err := s.CreateReport(rpt); err != nil {
		t.Error(err)
	}
	rpt.Other = "123"
	if err := s.UpdateReport(rpt); err != nil {
		t.Error(err)
	}

	get, err := s.GetReport(rpt.ID)
	if err != nil {
		t.Error(err)
	}
	if get.Other != rpt.Other {
		t.Error("Report update err, get:", get, "expected:", rpt)
	}
}

func beforeReport() *datasource {
	s := newTest()
	s.Exec("DELETE FROM reports")
	return s
}
