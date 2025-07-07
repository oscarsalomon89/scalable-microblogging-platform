package user

import (
	"context"
	"errors"
	"time"
)

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrUsernameExists     = errors.New("username already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrFolloweeNotFound   = errors.New("followee not found")
	ErrAlreadyFollowing   = errors.New("already following")
	ErrCannotFollowSelf   = errors.New("cannot follow self")
	ErrCannotUnfollowSelf = errors.New("cannot unfollow self")
	ErrNotFollowing       = errors.New("not following")
)

type (
	User struct {
		ID        string
		Username  string
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	//go:generate mockery --name=UserFinder --output=mocks --outpkg=mocks --filename=user_finder.go
	UserFinder interface {
		ExistsByUsername(ctx context.Context, username string) (bool, error)
		ExistsByID(ctx context.Context, id string) (bool, error)
		IsFollowing(ctx context.Context, followerID, followeeID string) (bool, error)
	}

	//go:generate mockery --name=UserCreator --output=mocks --outpkg=mocks --filename=user_creator.go
	UserCreator interface {
		CreateUser(ctx context.Context, user *User) error
		FollowUser(ctx context.Context, followerID, followeeID string) error
		UnfollowUser(ctx context.Context, followerID, followeeID string) error
	}

	//go:generate mockery --name=TimelineCache --output=mocks --outpkg=mocks --filename=timeline_cache_mock.go
	TimelineCache interface {
		InvalidateTimeline(ctx context.Context, userID string) error
	}
)
