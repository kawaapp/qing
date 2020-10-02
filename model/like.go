package model

const (
	LikePost = "post"
	LikeDiscussion = "dz"
)

type Like struct {
	ID        int64 `json:"id" meddler:"id,pk"`
	CreatedAt int64 `json:"created_at" meddler:"created_at"`

	// Status is favor state, 0 not 1 like
	Status int64 `json:"status" meddler:"status"`

	// Who clicks the like button
	UserID int64 `json:"user_id" meddler:"user_id"`

	// Like type, post or discussion
	TargetTy string `json:"target_ty" meddler:"target_ty"`

	// Like target's id
	TargetID int64 `json:"target_id" meddler:"target_id"`
}
