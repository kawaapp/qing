package datasource

import (
	"github.com/kawaapp/kawaqing/model"
	"github.com/russross/meddler"
	"fmt"
)

func (db *datasource) CreateMedia(m *model.Media) error {
	m.CreatedAt = UnixNow()
	return meddler.Insert(db, "medias", m)
}

func (db *datasource) GetMediaListByPostIds(pids []int64) ([]*model.Media, error) {
	q := joinIntArray(pids)
	data := make([]*model.Media, 0)
	err := meddler.QueryAll(db, &data, fmt.Sprintf(sqlListMediaByPostId, q))
	return data, err
}

func (db *datasource) GetMediaByPostId(pid int64) (*model.Media, error) {
	media := new(model.Media)
	err := meddler.QueryRow(db, media, sqlGetMediaByPostId, pid)
	return media, err
}

func (db *datasource) DeleteMediaByPostId(pid int64) error {
	stmt := sqlDeleteMediaByPostId
	_, err := db.Exec(stmt, pid)
	return err
}

const sqlGetMediaByPostId = `
SELECT
	id,
	created_at,
	post_id,
	author_id,
	_type,
	path,
	meta
FROM medias WHERE post_id=?
;`
const sqlDeleteMediaByPostId = `
DELETE FROM medias WHERE post_id=?
;`

const sqlListMediaByPostId = `
SELECT
	id,
	created_at,
	post_id,
	author_id,
	_type,
	path,
	meta
FROM medias
WHERE post_id IN (%s)
;`


