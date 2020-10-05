package datasource

import (
	"github.com/kawaapp/kawaqing/model"
	"github.com/russross/meddler"
	"fmt"
)

func (db *datasource) GetFavoriteList(q model.QueryParams, page, size int) ([]*model.Favorite, error)  {
	query, args := sqlGetFavoriteList(sqlFavoriteSelect, q, page, size)
	data := make([]*model.Favorite, 0)
	err := meddler.QueryAll(db, &data, query, args...)
	return data, err
}

func (db *datasource) GetFavoriteCount(q model.QueryParams) (int, error) {
	query, args := sqlGetFavoriteList(sqlFavoriteSelect, q, 0, 0)
	return Count(db, query, args...)
}

func (db *datasource) CreateFavorite(f *model.Favorite) error  {
	f.CreatedAt = UnixNow()
	return meddler.Insert(db, "favorites", f)
}

func (db *datasource) GetFavoriteUser(uid, pid int64) (*model.Favorite, error) {
	f := new(model.Favorite)
	err := meddler.QueryRow(db, f, sqlGetFavorite, uid, pid)
	return f, err
}

func (db *datasource) GetFavoriteId(id int64) (*model.Favorite, error)  {
	f := new(model.Favorite)
	err := meddler.Load(db, "favorites", f, id)
	return f, err
}

func (db *datasource) DeleteFavorite(id int64) error {
	_, err := db.Exec(sqlDeleteFavoriteId, id)
	return err
}

const sqlFavoriteSelect = `
SELECT
	id,
	created_at,
	user_id,
	discussion_id
`

const sqlGetFavorite = `
SELECT
	id,
	created_at,
	user_id,
	discussion_id
FROM favorites
WHERE user_id=? AND discussion_id=? LIMIT 1
`

func sqlGetFavoriteList(baseQuery string, params model.QueryParams, page, size int) (query string, args []interface{})  {
	query += baseQuery
	query += " FROM favorites"

	where := ""
	if q, ok := params["user_id"]; ok {
		where += " AND user_id=?"
		args = append(args, q)
	}
	if q, ok := params["discussion_id"]; ok {
		where += " AND discussion_id=?"
		args = append(args, q)
	}
	if len(where) > 0 {
		query += " WHERE 1=1 " + where
	}
	if size > 0 {
		query += fmt.Sprintf(" ORDER BY id DESC LIMIT %d OFFSET %d", size, page*size )
	}
	return
}

const sqlDeleteFavoriteId  = `
DELETE FROM favorites WHERE id=?
`