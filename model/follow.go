package model

type Follow struct {
	ID        int64 `json:"id" meddler:"id,pk"`
	CreatedAt int64 `json:"created_at" meddler:"created_at"`

	UserId     int64 `json:"user_id" meddler:"user_id"`
	FollowerId int64 `json:"follower_id" meddler:"follower_id"`
}