package model

import (
	"errors"
	"regexp"
)

// validate a username (e.g. from github)
var reUsername = regexp.MustCompile("^[a-zA-Z0-9-_.]+$")

var errUserLoginInvalid = errors.New("Invalid User Login")

// User represents a registered user.
type User struct {
	ID        int64 `json:"id" meddler:"id,pk"`
	CreatedAt int64 `json:"created_at" meddler:"created_at"`
	UpdatedAt int64 `json:"updated_at" meddler:"updated_at"`

	// counter cache
	SignCount int `json:"sign_count" meddler:"sign_count"`
	ExpCount  int `json:"exp_count"  meddler:"exp_count"`

	// Add +
	LastLogin int64 `json:"-" meddler:"last_login"`

	// Login is the username.
	Login    string `json:"login" meddler:"login"`
	Nickname string `json:"nickname" meddler:"nickname"`
	Email    string `json:"email" meddler:"email"`
	Phone    string `json:"phone" meddler:"phone"`
	Avatar   string `json:"avatar" meddler:"avatar"`
	Summary  string `json:"summary" meddler:"summary"`

	// User status
	BlockedAt  int64 `json:"blocked_at" meddler:"blocked_at"`
	SilencedAt int64 `json:"silenced_at" meddler:"silenced_at"`

	// Hash is a unique token used to sign tokens.
	Hash         string `json:"-" meddler:"hash"`
	PasswordHash string `json:"-" meddler:"password_hash"`
}

func (u *User) Validate() error {
	switch {
	case len(u.Login) == 0:
		return errUserLoginInvalid
	case len(u.Login) > 250:
		return errUserLoginInvalid
	case !reUsername.MatchString(u.Login):
		return errUserLoginInvalid
	default:
		return nil
	}
}

func (u *User) Blocked() bool {
	return u.BlockedAt > 0
}

func (u *User) Silenced() bool {
	return u.SilencedAt > 0
}

// UserBind represent the relationships between user and third-party Auth system.
type UserBind struct {
	ID        int64 `json:"id" meddler:"id,pk"`
	CreatedAt int64 `json:"-" meddler:"created_at"`

	Kind    string `json:"-" meddler:"kind"`
	UserId  int64  `json:"-" meddler:"user_id"`
	BindId  string `json:"-" meddler:"bind_id"`
	UnionId string `json:"-" meddler:"union_id"`
}
