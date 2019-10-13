package model

type Like struct {
	ID        int64 `json:"id" meddler:"id,pk"`
	CreatedAt int64 `json:"created_at" meddler:"created_at"`

	// Status is favor state, 0 not 1 like
	Status int64 `json:"status" meddler:"status"`

	// Who clicks the like button
	AuthorID int64 `json:"-" meddler:"author_id"`

	// Post is liked
	PostId int64 `json:"-" meddler:"post_id"`
}
