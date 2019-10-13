package model

type Category struct {
}

type Tag struct {
	ID        int64 `json:"id" meddler:"id,pk"`
	CreatedAt int64 `json:"created_at" meddler:"created_at"`

	// 用来排序
	Order   int    `json:"order" meddler:"_order"`

	// rgba
	Color   int    `json:"color" meddler:"color"`

	// 父节点
	ParentId int64 `json:"parent_id" meddler:"parent_id"`

	Text    string `json:"text" meddler:"text"`
	Summary string `json:"summary" meddler:"summary"`
}

// 关系表
type TagDiscussion struct {
	ID        int64 `json:"id" meddler:"id,pk"`
	CreatedAt int64 `json:"created_at" meddler:"created_at"`

	// relation
	DiscussionID int64 `json:"-" meddler:"discussion_id"`
	TagID  int64 `json:"-" meddler:"tag_id"`
}
