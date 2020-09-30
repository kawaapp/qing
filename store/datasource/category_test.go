package datasource

import (
	"testing"
	"github.com/kawaapp/kawaqing/model"
	"database/sql"
)

func TestCreateCategory(t *testing.T)  {
	s := beforeCategory()
	defer s.Close()

	ctg := &model.Category{
		Text: "分类1",
	}
	if err := s.CreateCategory(ctg); err != nil {
		t.Fatal(err)
	}
	get, err := s.GetCategory(ctg.ID)
	if err != nil {
		t.Fatal(err)
	}
	if get.Text != ctg.Text {
		t.Fatal(err)
	}

	// update
	ctg.Text = "分类2"
	if err := s.UpdateCategory(ctg); err != nil {
		t.Fatal(err)
	}
	if get, _ := s.GetCategory(ctg.ID); get.Text != ctg.Text {
		t.Fatal("UpdateCategory error, expect:", ctg, "get:", get)
	}

	// delete
	if err := s.DeleteCategory(ctg.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := s.GetCategory(ctg.ID); err != sql.ErrNoRows {
		t.Fatal("DeleteCategory error!")
	}
}

func beforeCategory() *datasource {
	s := newTest()
	s.Exec("DELETE FROM categories")
	return s
}
