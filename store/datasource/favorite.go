package datasource

import (
	"github.com/kawaapp/kawaqing/model"
	"github.com/russross/meddler"
)

func (db *datasource) GetFavoriteListUser(uid int64, page, size int) ([]*model.Favorite, error)  {
	if page == 0 {
		page = 1
	}
	data := make([]*model.Favorite, 0)
	err := meddler.QueryAll(db, &data, sqlListFavoriteUser, uid, size, (page-1) * size)
	return data, err
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

const sqlListFavoriteUser = `
SELECT
	id,
	created_at,
	user_id,
	discussion_id
FROM favorites
WHERE user_id = ?
ORDER BY id DESC LIMIT ? OFFSET ?
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

const sqlDeleteFavoriteId  = `
DELETE FROM favorites WHERE id=?
`