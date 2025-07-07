package tweet

import (
	"context"
	"time"
)

type (
	Tweet struct {
		ID        string
		UserID    string
		Content   string
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	//go:generate mockery --name=UserFinder --output=mocks --outpkg=mocks --filename=user_finder.go
	UserFinder interface {
		ExistsByID(ctx context.Context, id string) (bool, error)
		GetFollowers(ctx context.Context, id string) ([]string, error)
		GetFollowees(ctx context.Context, userID string) ([]string, error)
	}

	//go:generate mockery --name=TweetCreator --output=mocks --outpkg=mocks --filename=tweet_creator.go
	TweetCreator interface {
		CreateTweet(ctx context.Context, tweet *Tweet) error
	}

	//go:generate mockery --name=TweetReader --output=mocks --outpkg=mocks --filename=tweet_reader.go
	TweetReader interface {
		GetTweetsByUserIDs(ctx context.Context, userIDs []string, limit, offset int) ([]Tweet, error)
	}

	//go:generate mockery --name=TimelineCache --output=mocks --outpkg=mocks --filename=timeline_cache_mock.go
	TimelineCache interface {
		InvalidateTimeline(ctx context.Context, userID string) error
		GetTimeline(ctx context.Context, userID string) ([]Tweet, error)
		SetTimeline(ctx context.Context, userID string, tweets []Tweet) error
	}
)
