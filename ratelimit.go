package ratelimit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Config struct {
	MaxTokens int64
	Rate      int64
	Duration  time.Duration
	Redis     *redis.Client
}

func New(ctx context.Context, conf Config) RateLimit {
	return RateLimit{
		maxTokens:          conf.MaxTokens,
		lastTokenTimestamp: time.Now(),
		rds:                conf.Redis,
		duration:           conf.Duration,
	}
}

type RateLimit struct {
	rds                *redis.Client
	maxTokens          int64
	currentTokens      int64
	lastTokenTimestamp time.Time
	duration           time.Duration
}

func (b *RateLimit) GetBucket(ctx context.Context, identifier string) (Bucket, error) {
	bucket := Bucket{
		rds:                b.rds,
		MaxTokens:          b.maxTokens,
		CurrentTokens:      b.currentTokens,
		LastTokenTimestamp: b.lastTokenTimestamp,
		Identifier:         identifier,
		Duration:           b.duration,
	}

	brds, err := b.rds.Get(ctx, identifier).Result()
	if err != nil && err == redis.Nil {
		bts, _ := json.Marshal(bucket)
		if err := b.rds.Set(ctx, identifier, bts, b.duration).Err(); err != nil {
			return bucket, errors.New("invalid response from redis: " + err.Error())
		}
		return bucket, nil
	}

	if err := json.Unmarshal([]byte(brds), &bucket); err != nil {
		fmt.Println("invalid response from redis")
		return bucket, errors.New("invalid response from redis")
	}

	return bucket, nil
}
