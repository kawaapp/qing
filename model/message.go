package model

type NotType int

const (
	NotUnknown NotType = iota
	NotComment
	NotLike
)

type ChatType int

const (
	ChatText = iota
	ChatImage
	ChatAudio
	ChatVideo
)

// Forum notification
type Notification struct {
	ID        int64 `json:"id" meddler:"id,pk"`
	CreatedAt int64 `json:"created_at" meddler:"created_at"`

	// notification trigger
	EntityId int64 `json:"entity_id" meddler:"entity_id"`
	EntityType int `json:"entity_ty" meddler:"entity_ty"`

	// who send the message
	FromId int64 `json:"-" meddler:"from_id"`
	ToId   int64 `json:"-" meddler:"to_id"`

	// status 0 unread, 1 read
	Status int   `json:"status" meddler:"status"`
}

// Chat messages
type Chat struct {
	ID        int64 `json:"id" meddler:"id,pk"`
	CreatedAt int64 `json:"created_at" meddler:"created_at"`

	// message
	Content     string `json:"content" meddler:"content"`
	ContentType int64  `json:"type"    meddler:"_type"`

	// status
	Status int64 `json:"status" meddler:"status"`

	// 会话Id = A << 32 + B
	// A < B
	ChatId int64 `json:"-" meddler:"chat_id"`

	// who send the message
	FromId int64 `json:"-" meddler:"from_id"`
	ToId   int64 `json:"-" meddler:"to_id"`
}

type MessageCount struct {
	Favors int   `json:"favors" meddler:"favors"`
	Comments int `json:"comments" meddler:"comments"`
}