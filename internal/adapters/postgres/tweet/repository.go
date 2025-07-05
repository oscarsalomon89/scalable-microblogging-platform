package tweet

import (
	"context"
	"fmt"

	"github.com/oscarsalomon89/go-hexagonal/internal/application/tweet"
	db "github.com/oscarsalomon89/go-hexagonal/internal/platform/pg"
)

type tweetRepository struct {
	db db.Connections
}

func NewTweetRepository(db db.Connections) *tweetRepository {
	return &tweetRepository{db: db}
}

func (r *tweetRepository) CreateTweet(ctx context.Context, tweet *tweet.Tweet) error {
	tweetModel := fromDomain(tweet)

	err := r.db.MasterConn.
		WithContext(ctx).
		Create(tweetModel).Error
	if err != nil {
		return fmt.Errorf("failed to create tweet: %w", err)
	}

	tweet.ID = tweetModel.ID.String()
	tweet.CreatedAt = tweetModel.CreatedAt
	tweet.UpdatedAt = tweetModel.UpdatedAt

	return nil
}

func (r *tweetRepository) GetTweetsByUserIDs(ctx context.Context, userIDs []string, limit, offset int) ([]tweet.Tweet, error) {
	var tweets []Tweet
	if err := r.db.MasterConn.
		WithContext(ctx).
		Where("user_id IN ?", userIDs).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&tweets).Error; err != nil {
		return nil, fmt.Errorf("failed to get tweets: %w", err)
	}

	var tweetList []tweet.Tweet
	for _, tweet := range tweets {
		tweetList = append(tweetList, tweet.toDomain())
	}

	return tweetList, nil
}
