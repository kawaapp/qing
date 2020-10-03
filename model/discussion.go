package model

// 文章, 页面
type Discussion struct {
	ID        int64 `json:"id"         meddler:"id,pk"`
	CreatedAt int64 `json:"created_at" meddler:"created_at"`
	UpdatedAt int64 `json:"updated_at" meddler:"updated_at"`

	Title    string `json:"title"   meddler:"title"`
	Content  string `json:"content" meddler:"content"`
	AuthorID int64  `json:"author_id" meddler:"author_id"`

	// 状态
	Status int  `json:"status" meddler:"status"`

	// category
	CategoryID int64 `json:"cate_id" meddler:"cate_id"`

	// 最后回复人
	LastReplyUid int64 `json:"last_reply_uid" meddler:"last_reply_uid"`

	// 最后回复时间
	LastReplyAt int64 `json:"last_reply_at" meddler:"last_reply_at"`

	// 计数缓存
	CommentCount int   `json:"comment_count" meddler:"comment_count"`
	LikeCount    int   `json:"like_count" meddler:"like_count"`
	ViewCount    int   `json:"view_count" meddler:"view_count"`
}
