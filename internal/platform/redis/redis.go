package pkgredis

import (
	"context"
	"fmt"
	"time"

	"github.com/oscarsalomon89/go-hexagonal/internal/platform/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisConnection(cfg config.Cache) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return rdb, nil
}
