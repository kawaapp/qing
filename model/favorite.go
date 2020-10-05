package model

// 收藏
type Favorite struct {
	ID        int64 `json:"id" meddler:"id,pk"`
	CreatedAt int64 `json:"created_at" meddler:"created_at"`

	// Who clicks the like button
	UserID int64 `json:"user_id" meddler:"user_id"`

	// Entity is the id of post or comment
	DiscussionID int64 `json:"discussion_id" meddler:"discussion_id"`
}
