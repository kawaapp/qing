package datasource

import (
	"github.com/russross/meddler"
	"github.com/kawaapp/kawaqing/model"
	"fmt"
)

// post

func (db *datasource) GetPostList(params model.QueryParams, page, size int) ([]*model.Post, error) {
	data := make([]*model.Post, 0)
	query, args := sqlPostQuery("SELECT * ", params, page, size)
	err := meddler.QueryAll(db, &data, query, args...)
	return data, err
}

func (db *datasource) GetPostCount(params model.QueryParams) (int, error) {
	query, args := sqlPostQuery("SELECT COUNT(*) ", params, 0, 0)
	return db.Count(query, args...)
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

func sqlPostQuery(queryBase string, params model.QueryParams, page, size int) (query string, args []interface{}) {
	query += queryBase
	query += " FROM posts"

	where := ""
	if q, ok := params["content"]; ok {
		where += " AND content LIKE ?"
		args = append(args, "%" + q + "%")
	}

	if q, ok := params["author"]; ok {
		where += " AND author_id IN (SELECT id FROM users WHERE nickname=?)"
		args = append(args, q)
	}

	if len(where) > 0 {
		query += " WHERE 1=1" + where
	}

	if size > 0 {
		query += fmt.Sprintf(" ORDER BY id DESC LIMIT %d OFFSET %d", size, page * size)
	}
	return
}

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