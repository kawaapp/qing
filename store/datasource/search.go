package datasource

import (
	"github.com/kawaapp/kawaqing/model"
	"github.com/russross/meddler"
	"fmt"
)

func (db *datasource) SearchUser(params model.QueryParams, page, size int) ([]*model.User, error) {
	users := make([]*model.User, 0)
	query, args := sqlUserQuery("SELECT * ", params, page, size)
	err := meddler.QueryAll(db, &users, query, args...)
	return users, err
}

func (db *datasource) SearchUserCount(params model.QueryParams) (int, error) {
	query, args := sqlUserQuery("SELECT COUNT(id)", params, 0, 0)
	num, err := db.Count(query, args...)
	return num, err
}

func (db *datasource) SearchDiscussion(params model.QueryParams, page, size int) ([]*model.Discussion, error) {
	data := make([]*model.Discussion, 0)
	query, args := sqlDiscussionQuery("SELECT * ", params, page, size)
	err := meddler.QueryAll(db, &data, query, args...)
	return data, err
}

func (db *datasource) SearchDiscussionCount(params model.QueryParams) (int, error)  {
	query, args := sqlDiscussionQuery("SELECT COUNT(id) ", params, 0, 0)
	num, err := db.Count(query, args...)
	return num, err
}

func (db *datasource) SearchPost(params model.QueryParams, page, size int) ([]*model.Post, error) {
	data := make([]*model.Post, 0)
	query, args := sqlPostQuery("SELECT * ", params, page, size)
	err := meddler.QueryAll(db, &data, query, args...)
	return data, err
}

func (db *datasource) SearchPostCount(params model.QueryParams) (int, error) {
	query, args := sqlPostQuery("SELECT COUNT(*) ", params, 0, 0)
	return db.Count(query, args...)
}

func (db *datasource) SearchReport(params model.QueryParams, page, size int) ([]*model.Report, error) {
	reports := make([]*model.Report, 0)
	query, args := sqlReportQuery(sqlReportBase, params, page, size)
	err := meddler.QueryAll(db, &reports, query, args...)
	return reports, err
}

func (db *datasource) SearchReportCount(params model.QueryParams) (int, error) {
	query, args := sqlReportQuery(sqlReportCount, params, 0, 0)
	return db.Count(query, args...)
}

func (db *datasource) SearchSignUser(page, size int) ([]*model.User, error){
	data := make([]*model.User, 0)
	err := meddler.QueryAll(db, &data, sqlSearchUserSign, size, page * size)
	return data, err
}

func (db *datasource) SearchSignUserCount(page, size int) (int, error) {
	return  db.Count("SELECT COUNT(*) FROM users")
}

func (db *datasource) Count(stmt string, args ...interface{}) (int, error) {
	rows := db.QueryRow(stmt, args...)
	var count int
	err := rows.Scan(&count)
	return count, err
}

// user
func sqlUserQuery(queryBase string, params model.QueryParams, page, size int) (query string, args []interface{}) {
	query += queryBase
	query += " FROM users"

	where := ""
	if q, ok := params["login"]; ok {
		where += " AND login=?"
		args = append(args, q)
	}

	if 	q, ok := params["nickname"]; ok {
		where += " AND nickname LIKE ?"
		args = append(args, "%" + q + "%")
	}

	if _, ok := params["silence"]; ok {
		where += " AND silenced_at > ?"
	}

	if _, ok := params["block"]; ok {
		where += " AND blocked_at > 0"
	}

	if len(where) > 0 {
		query += " WHERE 1=1" + where
	}

	if size > 0 {
		query += fmt.Sprintf(" ORDER BY id DESC LIMIT %d OFFSET %d", size, page * size)
	}
	return
}

func sqlDiscussionQuery(queryBase string, params model.QueryParams, page, size int) (query string, args []interface{}) {
	query += queryBase
	query += " FROM discussions"

	where := ""
	if q, ok := params["content"]; ok {
		where += " AND content LIKE ?"
		args = append(args, "%" + q + "%")
	}

	// TODO JOIN!
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

const sqlReportBase = `
SELECT
	c.id,
	c.created_at,
	c.updated_at,
	c.entity_id,
	c.entity_ty,
	c.content,
	c.counter,
	c.status,
	c.user_id,
	c.report_ty,
	c.other,
	c.images
`

const sqlReportCount = `
SELECT COUNT(c.id)
`

func sqlReportQuery(queryBase string, params model.QueryParams, page, size int) (query string, args []interface{}) {
	query += queryBase
	query += " FROM reports c"

	// join
	query += " LEFT JOIN posts a ON a.id = c.entity_id"
	query += " LEFT JOIN users b ON b.id = c.user_id"

	// where
	where := ""
	if q, ok := params["post"]; ok {
		where += " AND a.content LIKE ?"
		args = append(args, "%" + q + "%")
	}
	if q, ok := params["user"]; ok {
		where += " AND b.nickname LIKE ?"
		args = append(args, "%" + q + "%")
	}
	if _, ok := params["status"]; ok {
		where += " AND c.status > 0"
	}

	if len(where) > 0 {
		query += " WHERE 1=1" + where
	}

	if size > 0 {
		query += fmt.Sprintf(" ORDER BY c.id DESC LIMIT %d OFFSET %d", size, page * size)
	}
	return
}

const sqlSearchUserSign = `
SELECT
	users.id,
	users.created_at,
	login,
	nickname,
	avatar,
	sign_count
FROM users
ORDER BY sign_count DESC, id DESC LIMIT ? OFFSET ?
`