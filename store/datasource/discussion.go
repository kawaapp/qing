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
	num, err := db.Count(query, args...)
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
	comment_count
`

const sqlDiscussionList = sqlDiscussionSelect +`
FROM discussions
ORDER BY id DESC LIMIT ? OFFSET ?;`

const sqlListDiscussionComments = sqlDiscussionSelect + `
FROM discussions
ORDER BY comment_count DESC, id DESC LIMIT ? OFFSET ?
;`

const sqlListDiscussionByUser = sqlDiscussionSelect +`
FROM discussions
WHERE author_id = ?
ORDER BY id DESC LIMIT ? OFFSET ?
;`

const sqlListPostByIds = sqlDiscussionSelect +`
FROM discussions
WHERE id IN (%s) ORDER BY id DESC
;`

const sqlDiscussionDelete = `DELETE FROM discussions WHERE id=?;`