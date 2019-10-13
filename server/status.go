package server

// Account
const (
	AccountUnknown = 1000 + iota
	AccountBlocked
	AccountSilenced
)

// Post
const (
	PostUnknown = 2000 + iota
	PostInvalid
)

// Comment
const (
	CommentUnknown = 3000 + iota
	CommentInvalid
)

// Sign-in
const (
	SignUnknown = 4000 + iota
	SignRepeat
)