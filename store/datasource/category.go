package datasource

import (
	"github.com/kawaapp/kawaqing/model"
	"github.com/russross/meddler"
)

func (db *datasource) GetCategoryList() ([]*model.Category, error) {
	return nil, nil
}

func (db *datasource) GetCategory(id int64) (*model.Category, error) {
	ctg := new(model.Category)
	err := meddler.Load(db,"categories", ctg, id)
	return ctg, err
}

func (db *datasource) CreateCategory(ctg *model.Category) error {
	ctg.CreatedAt = UnixNow()
	ctg.UpdatedAt = UnixNow()
	return meddler.Insert(db,"categories", ctg)
}

func (db *datasource) UpdateCategory(ctg *model.Category) error {
	ctg.UpdatedAt = UnixNow()
	return meddler.Update(db, "categories", ctg)
}

func (db *datasource) DeleteCategory(id int64) error {
	return Delete(db,"DELETE FROM categories WHERE id=?", id)
}