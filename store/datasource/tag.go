package datasource

import (
	"github.com/russross/meddler"
	"github.com/kawaapp/kawaqing/model"
)

// find tag by name
func (db *datasource) GetTag(tag string) (*model.Tag, error) {
	stmt := sqlTagFindByName
	t := &model.Tag{}
	err := meddler.QueryRow(db, t, stmt, tag)
	return t, err
}

func (db *datasource) GetTagId(id int64) (*model.Tag, error) {
	t := new(model.Tag)
	err := meddler.Load(db, "tags", t, id)
	return t, err
}

// get or insert tag
func (db *datasource) CreateTag(tag, summary string) (*model.Tag, error) {
	t, err := db.GetTag(tag)
	if err != nil {
		m := &model.Tag{Text: tag, Summary:summary, CreatedAt: UnixNow()}
		if err := meddler.Insert(db, "tags", m); err != nil {
			return nil, err
		}
		t = m
	}
	return t, nil
}

func (db *datasource) DeleteTag(id int64) error {
	_, err := db.Exec(sqlDeleteTag, id)
	return err
}

func (db *datasource) UpdateTag(t *model.Tag) error {
	return meddler.Update(db, "tags", t)
}

func (db *datasource) GetDiscussionsByTag(tag string, page, size int) ([]*model.Discussion, error) {
	t, err := db.GetTag(tag)
	if err != nil {
		return nil, err
	}

	// get posts
	stmt := sqlTagListDiscussion
	posts := make([]*model.Discussion, 0)
	err = meddler.QueryAll(db, &posts, stmt, t.ID, page, size)
	return posts, err
}

func (db *datasource) GetTagList() ([]*model.Tag, error) {
	stmt := sqlTagList
	tags := make([]*model.Tag, 0)
	err := meddler.QueryAll(db, &tags, stmt)
	return tags, err
}

func (db *datasource) LinkTagDiscussion(pid int64, tags []string) error {
	for _, t := range tags {
		if err := db.linkTagPost(pid, t); err != nil {
			return err
		}
	}
	return nil
}

// make tag.text unique string
func (db *datasource) linkTagPost(pid int64, tag string) error {
	t, err := db.CreateTag(tag, "")
	if err != nil {
		return err
	}

	// create relations
	rel := model.TagDiscussion{
		DiscussionID: pid,
		TagID:  t.ID,
		CreatedAt:UnixNow(),
	}
	return meddler.Insert(db, "tag_discussions", &rel)
}

const sqlDeleteTag = `
DELETE FROM tags WHERE id=?
;`


const sqlTagFindByName = `
SELECT
	id,
	created_at,
	_order,
	text,
	color,
	summary
FROM tags
WHERE text=?;`

const sqlTagListDiscussion = sqlDiscussionSelect + `
FROM discussions
WHERE id
IN(
	SELECT discussion_id FROM tag_discussions WHERE tag_id=? AND discussion_id<?
) ORDER BY id DESC LIMIT ?;`


const sqlTagList = `
SELECT
	id,
	created_at,
	_order,
	text,
	summary
FROM tags ORDER BY _order DESC;`