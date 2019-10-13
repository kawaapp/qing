package model

//
type Post struct {
	ID        int64 `json:"id" meddler:"id,pk"`
	CreatedAt int64 `json:"created_at" meddler:"created_at"`

	// 讨论ID
	DiscussionID   int64  `json:"discussion_id" meddler:"discussion_id"`

	// 被回复的评论ID
	ParentID int64  `json:"parent_id" meddler:"parent_id"`

	// 回复的人
	AuthorID int64  `json:"author_id" meddler:"author_id"`

	// 被回复的人
	ReplyID  int64  `json:"reply_id" meddler:"reply_id"`

	// 点赞数
	LikeCount int `json:"-" meddler:"like_count"`

	// 评论内容
	Content  string `json:"content" meddler:"content"`
}
