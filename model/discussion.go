package model

// 文章, 页面
type Discussion struct {
	ID        int64 `json:"id"         meddler:"id,pk"`
	CreatedAt int64 `json:"created_at" meddler:"created_at"`
	UpdatedAt int64 `json:"updated_at" meddler:"updated_at"`

	Title    string `json:"title"   meddler:"title"`
	Content  string `json:"content" meddler:"content"`
	AuthorID int64  `json:"author_id" meddler:"author_id"`

	// category
	CategoryID int64 `json:"cate_id" meddler:"cate_id"`

	Status int  `json:"status" meddler:"status"`

	// Add+
	FirstPost    int64 `json:"first_post" meddler:"first_post"`
	LastPost     int64 `json:"last_post" meddler:"last_post"`
	CommentCount int   `json:"comment_count" meddler:"comment_count"`

	LikeCount    int   `json:"like_count" meddler:"like_count"`
	ViewCount    int   `json:"view_count" meddler:"view_count"`
}
