package datasource

import (
	"github.com/russross/meddler"
	"github.com/kawaapp/kawaqing/model"
	"fmt"
	"database/sql"
)

// post
func (db *datasource) GetLikeList(q model.QueryParams, page, size int) ([]*model.Like, error) {
	query, args := sqlGetLikeList(sqlLikeSelect, q, page, size)
	data := make([]*model.Like, 0)
	meddler.QueryAll(db, &data, query, args...)
	return data, nil
}

func (db *datasource) GetLikeCount(q model.QueryParams) (int, error) {
	query, args := sqlGetLikeList(sqlLikeSelect, q, 0, 0)
	rows := db.QueryRow(query, args...)
	var count int
	err := rows.Scan(&count)
	return count, err
}

func (db *datasource) GetUserLikedPostList(uid int64, pids []int64) ([]int64, error) {
	data := make([]int64, len(pids))

	q := joinIntArray(pids)
	stmt := fmt.Sprintf(sqlUserLikedPostList, q)
	rows, err := db.Query(stmt, uid, model.LikePost)
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

func (db *datasource) GetLike(t string, tid, uid int64) (*model.Like, error) {
	like := new(model.Like)
	err := meddler.QueryRow(db, like, sqlGetLike, t, tid, uid)
	return like, err
}

func (db *datasource) GetLikeId(id int64) (*model.Like, error) {
	f := new(model.Like)
	err := meddler.Load(db, "likes", f, id)
	return f, err
}

// create favor
func (db *datasource) CreateLike(t string, tid, uid int64) (error, bool) {
	firstTime := false
	err := db.Transact(func(tx *sql.Tx) error {
		l := new(model.Like)
		err := meddler.QueryRow(tx, l, sqlGetLike, t, tid, uid)
		if err != sql.ErrNoRows && err != nil {
			return err
		}
		// update it
		if err == nil {
			l.Status = 1
			return meddler.Update(tx, "likes", l)
		}

		// create it
		firstTime = true
		l = &model.Like {
			CreatedAt: UnixNow(),
			Status: 1,
			UserID: uid,
			TargetTy: t,
			TargetID: tid,
		}
		return meddler.Insert(tx, "likes", l)
	})
	return err, firstTime
}

func (db *datasource) DeleteLike(t string, tid, uid int64) error {
	return db.Transact(func(tx *sql.Tx) error {
		l := new(model.Like)
		err := meddler.QueryRow(tx, l, sqlGetLike, t, tid, uid)
		if err == sql.ErrNoRows {
			return nil
		}
		if err != nil {
			return err
		}
		if l.Status == 0 {
			return nil
		}

		//set status = 0
		l.Status = 0
		return meddler.Update(tx, "likes", l)
	})
}

func (db *datasource) deleteLike(id int64) error {
	return Delete(db, `DELETE FROM likes WHERE id=?`, id)
}

const sqlLikeSelect = `
SELECT
	id,
	created_at,
	status,
	user_id,
	target_ty,
	target_id
`

const sqlGetLike = sqlLikeSelect + `
	FROM likes
	WHERE target_ty=? AND target_id=? AND user_id=? LIMIT 1
`

func sqlGetLikeList(queryBase string, q model.QueryParams, page, size int) (query string, args []interface{}) {
	query += queryBase
	query += "FROM likes"

	where := ""
	if q, ok := q["target_ty"]; ok {
		where += " AND target_ty=?"
		args = append(args, q)
	}
	if q, ok := q["target_id"]; ok {
		where += " AND target_id=?"
		args = append(args, q)
	}
	if q, ok := q["user_id"]; ok {
		where += " AND user_id=?"
		args = append(args, q)
	}
	if len(where) > 0 {
		query += " status=1 " + where
	}
	if size > 0 {
		query += fmt.Sprintf(" ORDER BY id DESC LIMIT %d OFFSET %d", size, page*size)
	}
	return
}

const sqlUserLikedPostList = `
SELECT target_id
FROM likes
WHERE target_id in(%s) AND user_id=? AND target_ty=? AND status=1
`