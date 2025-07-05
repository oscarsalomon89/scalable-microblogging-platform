package user

import (
	"errors"
	"time"
)

var (
	ErrInvalidInput     = errors.New("invalid input")
	ErrUsernameExists   = errors.New("username already exists")
	ErrUserNotFound     = errors.New("user not found")
	ErrFolloweeNotFound = errors.New("followee not found")
	ErrAlreadyFollowing = errors.New("already following")
	ErrCannotFollowSelf = errors.New("cannot follow self")
)

type User struct {
	ID        string
	Username  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
