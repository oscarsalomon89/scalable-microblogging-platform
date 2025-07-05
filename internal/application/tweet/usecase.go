package tweet

import (
	"context"
	"fmt"

	"github.com/oscarsalomon89/go-hexagonal/internal/application/user"
	twcontext "github.com/oscarsalomon89/go-hexagonal/pkg/context"
)

type (
	UsersFinder interface {
		ExistsByID(ctx context.Context, id string) (bool, error)
		GetFollowers(ctx context.Context, id string) ([]string, error)
		GetFollowees(ctx context.Context, userID string) ([]string, error)
	}

	TweetsCreator interface {
		CreateTweet(ctx context.Context, tweet *Tweet) error
	}

	TweetReader interface {
		GetTweetsByUserIDs(ctx context.Context, userIDs []string, limit, offset int) ([]Tweet, error)
	}

	TimelineCache interface {
		InvalidateTimeline(ctx context.Context, userID string) error
		GetTimeline(ctx context.Context, userID string) ([]Tweet, error)
		SetTimeline(ctx context.Context, userID string, tweets []Tweet) error
	}

	usecase struct {
		userFinder    UsersFinder
		tweetReader   TweetReader
		tweetsCreator TweetsCreator
		cache         TimelineCache
	}
)

func NewTweetUseCase(userFinder UsersFinder, tweetReader TweetReader, tweetsCreator TweetsCreator, cache TimelineCache) *usecase {
	return &usecase{
		userFinder:    userFinder,
		tweetReader:   tweetReader,
		tweetsCreator: tweetsCreator,
		cache:         cache,
	}
}

func (uc *usecase) CreateTweet(ctx context.Context, tweet *Tweet) error {
	if exist, err := uc.userFinder.ExistsByID(ctx, tweet.UserID); err != nil {
		return fmt.Errorf("failed to check user ID: %w", err)
	} else if !exist {
		return user.ErrUserNotFound
	}

	if err := uc.tweetsCreator.CreateTweet(ctx, tweet); err != nil {
		return fmt.Errorf("failed to create tweet: %w", err)
	}

	detachedCtx := twcontext.NewDetachedWithRequestID(ctx)
	go uc.invalidateFollowersTimelinesAsync(detachedCtx, tweet.UserID)

	return nil
}

func (uc *usecase) GetTimeline(ctx context.Context, userID string, limit, offset int) ([]Tweet, error) {
	logger := twcontext.Logger(ctx)

	tweets, err := uc.cache.GetTimeline(ctx, userID)
	if err != nil {
		logger.WithError(err).Warn("failed to get timeline from cache")
	} else {
		if len(tweets) > 0 {
			return tweets, nil
		}
		logger.Info("timeline cache hit but empty")
		return []Tweet{}, nil
	}

	followeeIDs, err := uc.userFinder.GetFollowees(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get followees: %w", err)
	}
	if len(followeeIDs) == 0 {
		return []Tweet{}, nil
	}

	tweets, err = uc.tweetReader.GetTweetsByUserIDs(ctx, followeeIDs, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error retrieving timeline from Cassandra: %w", err)
	}

	if len(tweets) == 0 {
		return []Tweet{}, nil
	}

	if err := uc.cache.SetTimeline(ctx, userID, tweets); err != nil {
		logger.WithError(err).Error("Failed to set timeline cache")
	}

	return tweets, nil
}
