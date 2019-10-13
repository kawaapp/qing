package datasource

import (
	"github.com/russross/meddler"
	"github.com/kawaapp/kawaqing/model"
	"fmt"
)

// post
func (db *datasource) GetLikeList(pid int64, page, size int) ([]*model.Like, error) {
	stmt := sqlLikeList
	data := make([]*model.Like, 0)
	meddler.QueryAll(db, &data, stmt, pid, size, page * size)
	return data, nil
}

func (db *datasource) GetLikeListUser(uid int64, page, size int) ([]*model.Like, error) {
	stmt := sqlListLikeByUser
	data := make([]*model.Like, 0)
	err := meddler.QueryAll(db, &data, stmt, uid, size, page * size)
	return data, err
}

func (db *datasource) GetLikeCount(pid int64) (int, error) {
	stmt := sqlLikeCount
	rows := db.QueryRow(stmt, pid)
	var count int
	err := rows.Scan(&count)
	return count, err
}

func (db *datasource) GetLikePostList(uid int64, pids []int64) ([]int64, error) {
	data := make([]int64, len(pids))

	q := joinIntArray(pids)
	stmt := fmt.Sprintf(sqlLikedPostList, q)
	rows, err := db.Query(stmt, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	index := 0
	for rows.Next() {
		if err := rows.Scan(&data[index]); err != nil {
			return nil, err
		}
		index++
	}
	return data[:index], err
}

func (db *datasource) GetLike(pid, uid int64) (*model.Like, error) {
	favor := new(model.Like)
	err := meddler.QueryRow(db, favor, sqlGetLikeUserEntity, uid, pid)
	return favor, err
}

func (db *datasource) GetLikeId(id int64) (*model.Like, error) {
	f := new(model.Like)
	err := meddler.Load(db, "likes", f, id)
	return f, err
}

// create favor
func (db *datasource) CreateLike(f *model.Like) error {
	f.CreatedAt = UnixNow()
	return meddler.Insert(db, "likes", f)
}

func (db *datasource) UpdateLike(f *model.Like) error  {
	f.CreatedAt = UnixNow()
	return meddler.Update(db, "likes", f)
}

func (db *datasource) DeleteLike(id int64) error {
	stmt := sqlDeleteLike
	_, err := db.Exec(stmt, id)
	return err
}

const sqlDeleteLike = `DELETE FROM likes WHERE id=?`

const sqlLikeSelect = `
SELECT
	id,
	created_at,
	status,
	author_id,
	post_id
`

const sqlLikeList = sqlLikeSelect + `
FROM likes
WHERE status=1 AND post_id=?
ORDER BY id DESC LIMIT ? OFFSET ?;`


const sqlListLikeByUser = sqlLikeSelect + `
FROM likes
WHERE status=1 AND author_id=?
ORDER BY id DESC LIMIT ? OFFSET ?;`

const sqlLikeCount = `
SELECT
	COUNT(*)
FROM likes
WHERE status=1 AND post_id = ?
;`

const sqlGetLikeUserEntity = sqlLikeSelect + `
FROM likes
WHERE author_id=? AND post_id=? LIMIT 1
;`

const sqlLikedPostList = `
SELECT post_id
FROM likes
WHERE post_id in(%s) AND author_id=? AND status=1
`