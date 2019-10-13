package datasource

import (
	"github.com/russross/meddler"
	"github.com/kawaapp/kawaqing/model"
)

// 消息状态： 1 已读  0 未读

// Chat message
func (db *datasource) CreateChatMessage(m *model.Chat) error {
	return meddler.Insert(db, "chats", m)
}

func (db *datasource) GetChatMessageById(id int64) (*model.Chat, error) {
	chat := new(model.Chat)
	err := meddler.Load(db, "chats", chat, id)
	return chat, err
}

// 跟我有关的所有对话
func (db *datasource) GetChatUserList(to int64, page, size int) ([]*model.Chat, error) {
	stmt := sqlChatUserList
	data := make([]*model.Chat, 0)
	err := meddler.QueryAll(db, &data, stmt, to, size, page * size)
	return data, err
}

// 返回用户聊天信息
func (db *datasource) GetChatMsgList(from, to int64, page, size int) ([]*model.Chat, error) {
	var (
		stmt = sqlChatListFromTo
		chat = getChatId(from, to)
	)
	msgs := make([]*model.Chat, 0)
	err := meddler.QueryAll(db, &msgs, stmt, chat, size, page * size)
	return msgs, err
}


func (db *datasource) SetChatMsgAsRead(from, to int64) error {
	stmt := sqlChatReadFromTo
	_, err := db.Exec(stmt, from, to)
	return err
}

func (db *datasource) SetChatReadId(to, id int64) error {
	stmt := sqlChatReadId
	_, err := db.Exec(stmt, to, id)
	return err
}

func (db *datasource) DeleteChat(id int64, to int64) error {
	stmt := sqlDeleteChat
	_, err := db.Exec(stmt, id, to)
	return err
}

func getChatId(from, to int64) int64 {
	min, max := from, to
	if min > max {
		min, max = to, from
	}
	return min << 32 + max
}


const sqlChatSelect = `
SELECT
	id,
	created_at,
	content,
	_type,
	chat_id,
	status,
	from_id, to_id
`

const sqlChatUserList = sqlChatSelect + `
FROM chats
WHERE id IN
	(
		SELECT MAX(id) FROM chats WHERE ? IN (to_id, from_id) GROUP BY chat_id
	)
ORDER BY id DESC LIMIT ? OFFSET ?
;`

const sqlChatListFromTo = sqlChatSelect + `
FROM chats
WHERE chat_id=?
ORDER BY id DESC LIMIT ? OFFSET ?
;`

const sqlChatReadFromTo = `
UPDATE chats
SET
	status=1
WHERE from_id=? AND to_id=? AND status=0;`

const sqlChatReadId = `
UPDATE msg_box
SET
	status=1
WHERE to_id=? AND id=? AND status=0;`

const sqlDeleteChat = `DELETE FROM chats WHERE id=?`
