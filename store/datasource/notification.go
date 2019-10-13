package datasource

import (
	"github.com/russross/meddler"
	"github.com/kawaapp/kawaqing/model"
)

// 消息状态： 1 已读  0 未读

// Notification message
func (db *datasource) CreateNotification(m *model.Notification) error {
	m.CreatedAt = UnixNow()
	return meddler.Insert(db, "notifications", m)
}

func (db *datasource) GetNotificationById(id int64) (*model.Notification, error) {
	m := new(model.Notification)
	err := meddler.Load(db, "notifications", m, id)
	return m, err
}

func (db *datasource) GetNotificationCount(to int64) (*model.MessageCount, error) {
	stmt := sqlNotificationCount
	count := new(model.MessageCount)
	err := meddler.QueryRow(db, count, stmt, to, to)
	return count, err
}

// 返回用户的评论/点赞等系统消息
func (db *datasource) GetNotificationListType(to int64, mt model.NotType, page, size int) ([]*model.Notification, error) {
	stmt := sqlNotificationList
	msgs := make([]*model.Notification, 0)
	err := meddler.QueryAll(db, &msgs, stmt, mt, to, size, page * size)
	return msgs, err
}

func (db *datasource) SetNotificationReadId(uid, id int64) error {
	stmt := sqlNotificationReadId
	_, err := db.Exec(stmt, uid, id)
	return err
}


func (db *datasource) SetNotificationReadType(to int64, mt model.NotType) error {
	stmt := sqlNotificationReadType
	_, err := db.Exec(stmt, to, mt)
	return err
}

func (db *datasource) DeleteNotification(id int64) error {
	stmt := sqlDeleteNotification
	_, err := db.Exec(stmt, id)
	return err
}

const sqlNotificationCount = `
SELECT
	(
		SELECT COUNT(*)
		FROM notifications
		WHERE status=0 AND entity_ty=2 AND to_id=?
	) AS favors,
	(
		SELECT COUNT(*)
		FROM notifications
		WHERE status=0 AND entity_ty=1 AND to_id=?
	) AS comments
`

const sqlNotificationSelect = `
SELECT
	id,
	created_at,
	entity_id,
	entity_ty,
	from_id, to_id
`

const sqlNotificationList = sqlNotificationSelect + `
FROM notifications
WHERE entity_ty = ? AND to_id = ?
ORDER BY id DESC LIMIT ? OFFSET ?
;`

const sqlNotificationReadId = `
UPDATE notifications
SET
	status=1
WHERE to_id=? AND id=? AND status=0;`

const sqlNotificationReadType = `
UPDATE notifications
SET
	status=1
WHERE to_id=? AND status=0 AND entity_ty=?;`

const sqlDeleteNotification = `DELETE FROM notifications WHERE id=?`