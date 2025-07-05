package timeline

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/oscarsalomon89/go-hexagonal/internal/application/tweet"
	"github.com/redis/go-redis/v9"
)

type timelineCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewCache(c *redis.Client, ttl time.Duration) (*timelineCache, error) {
	return &timelineCache{client: c, ttl: ttl}, nil
}

func (r *timelineCache) GetTimeline(ctx context.Context, userID string) ([]tweet.Tweet, error) {
	key := fmt.Sprintf("timeline:%s", userID)

	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("timeline not found for user %s: %w", userID, err)
		}
		return nil, fmt.Errorf("failed to retrieve timeline cache for user %s: %w", userID, err)
	}

	var timelineCache []tweet.Tweet
	if err := json.Unmarshal([]byte(data), &timelineCache); err != nil {
		return nil, fmt.Errorf("failed to parse timeline data for user %s: %w", userID, err)
	}

	sort.Slice(timelineCache, func(i, j int) bool {
		return timelineCache[i].CreatedAt.After(timelineCache[j].CreatedAt)
	})

	return timelineCache, nil
}

func (r *timelineCache) SetTimeline(ctx context.Context, userID string, tweets []tweet.Tweet) error {
	key := fmt.Sprintf("timeline:%s", userID)

	data, err := json.Marshal(tweets)
	if err != nil {
		return fmt.Errorf("failed to serialize timeline data for user %s: %w", userID, err)
	}

	if err := r.client.Set(ctx, key, string(data), r.ttl).Err(); err != nil {
		return fmt.Errorf("failed to set timeline cache for user %s: %w", userID, err)
	}

	return nil
}

func (r *timelineCache) InvalidateTimeline(ctx context.Context, userID string) error {
	key := fmt.Sprintf("timeline:%s", userID)
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to invalidate timeline cache for user %s: %w", userID, err)
	}
	return nil
}
