package datasource

import (
	"github.com/russross/meddler"
	"github.com/kawaapp/kawaqing/model"
	"fmt"
)

// post
func (db *datasource) GetPostList(did int64, page, size int) ([]*model.Post, error) {
	stmt := sqlPostList
	data := make([]*model.Post, 0)
	err := meddler.QueryAll(db, &data, stmt, did, size, page * size)
	return data, err
}

func (db *datasource) GetPostListUser(uid int64, page, size int) ([]*model.Post, error) {
	stmt := sqlPostListByUser
	data := make([]*model.Post, 0)
	err := meddler.QueryAll(db, &data, stmt, uid, size, page * size)
	return data, err
}

func (db *datasource) GetPostListByIds(ids []int64) ([]*model.Post, error) {
	if len(ids) == 0 {
		return []*model.Post{}, nil
	}

	stmt := sqlPostListByIds
	q := joinIntArray(ids)
	data := make([]*model.Post, 0)
	err := meddler.QueryAll(db, &data, fmt.Sprintf(stmt, q))
	return data, err
}

func (db *datasource) GetPostCount(pid int64) (int, error) {
	stmt := sqlCommentCount
	rows := db.QueryRow(stmt, pid)
	var count int
	err := rows.Scan(&count)
	return count, err
}

func (db *datasource) GetPost(id int64) (*model.Post, error) {
	data := new(model.Post)
	err := meddler.Load(db, "posts", data, id)
	return data, err
}

func (db *datasource) CreatePost(p *model.Post) error {
	p.CreatedAt = UnixNow()
	return meddler.Insert(db, "posts", p)
}

func (db *datasource) UpdatePost(p *model.Post) error {
	return meddler.Update(db, "posts", p)
}

func (db *datasource) DeletePost(id int64) error {
	stmt := sqlPostDelete
	_, err := db.Exec(stmt, id)
	return err
}

// 这么看的话，直接缓存 Subject 好像方便多了...
func (db *datasource) GetCommentSubjectList(pids, cids []int64) ([]string, error) {
	return nil, nil
}

const sqlPostSelect = `
SELECT
	id,
	created_at,
	discussion_id,
	parent_id,
	author_id,
	reply_id,
	like_count,
	content
`

const sqlPostList = sqlPostSelect + `
FROM posts
WHERE discussion_id=?
ORDER BY id DESC LIMIT ? OFFSET ?
;`

const sqlPostListByUser = sqlPostSelect + `
FROM posts
WHERE author_id=?
ORDER BY id DESC LIMIT ? OFFSET ?
;`


const sqlPostListByIds = sqlPostSelect + `
FROM posts
WHERE id IN (%s);`

const sqlCommentCount = `
SELECT
	COUNT(*)
FROM posts WHERE post_id=?
;`

const sqlPostDelete = `DELETE FROM posts WHERE id=?;`