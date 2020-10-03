package datasource

import (
	"github.com/russross/meddler"
	"github.com/kawaapp/kawaqing/model"
	"fmt"
)

// posts
func (db *datasource) GetDiscussionList(params model.QueryParams, page, size int) ([]*model.Discussion, error) {
	data := make([]*model.Discussion, 0)
	query, args := sqlDiscussionQuery("SELECT * ", params, page, size)
	err := meddler.QueryAll(db, &data, query, args...)
	return data, err
}

func (db *datasource) GetDiscussionCount(params model.QueryParams) (int, error)  {
	query, args := sqlDiscussionQuery("SELECT COUNT(id) ", params, 0, 0)
	num, err := Count(db, query, args...)
	return num, err
}

func (db *datasource) GetDiscussionListByIds(ids []int64) ([]*model.Discussion, error) {
	if len(ids) == 0 {
		return []*model.Discussion{}, nil
	}
	stmt := fmt.Sprintf(sqlListPostByIds, joinIntArray(ids))
	data := make([]*model.Discussion,0)
	err := meddler.QueryAll(db, &data, stmt)
	return data, err
}

func (db *datasource) GetDiscussion(id int64) (*model.Discussion, error) {
	data := new(model.Discussion)
	err := meddler.Load(db, "discussions", data, id)
	return data, err
}

func (db *datasource) CreateDiscussion(p *model.Discussion) error {
	p.CreatedAt = UnixNow()
	p.UpdatedAt = UnixNow()
	return meddler.Insert(db, "discussions", p)
}

func (db *datasource) UpdateDiscussion(p *model.Discussion) error {
	return meddler.Update(db, "discussions", p)
}

func (db *datasource) DeleteDiscussion(id int64) error {
	stmt := sqlDiscussionDelete
	_, err := db.Exec(stmt, id)
	return err
}

const sqlDiscussionSelect = `
SELECT
	id,
	created_at,

	title,
	content,
	author_id,

	first_post,
	last_post,
	comment_count,
	view_count,
	like_count
`

func sqlDiscussionQuery(queryBase string, params model.QueryParams, page, size int) (query string, args []interface{}) {
	query += queryBase
	query += " FROM discussions"

	where := ""
	if q, ok := params["content"]; ok {
		where += " AND content LIKE ?"
		args = append(args, "%" + q + "%")
	}
	if q, ok := params["author_id"]; ok {
		where += " AND author_id = ?"
		args = append(args, q)
	}
	if q, ok := params["cate_id"]; ok {
		where += " AND cate_id = ?"
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

const sqlListPostByIds = sqlDiscussionSelect +`
FROM discussions
WHERE id IN (%s) ORDER BY id DESC
;`

const sqlDiscussionDelete = `DELETE FROM discussions WHERE id=?;`