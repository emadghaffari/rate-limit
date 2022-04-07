package ratelimit

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type Config struct {
	MaxTokens int64
	Rate      int64
	Duration  time.Duration
	Redis     *redis.Client
}

func New(ctx context.Context, conf Config) RateLimit {
	return RateLimit{
		rate:               conf.Rate,
		maxTokens:          conf.MaxTokens,
		lastTokenTimestamp: time.Now(),
		rds:                conf.Redis,
		duration:           conf.Duration,
	}
}

type RateLimit struct {
	rds                *redis.Client
	maxTokens          int64
	rate               int64
	currentTokens      int64
	lastTokenTimestamp time.Time
	duration           time.Duration
}

func (b *RateLimit) GetBucket(ctx context.Context, identifier string) Bucket {
	bucket := Bucket{
		rds:                b.rds,
		MaxTokens:          b.maxTokens,
		Rate:               b.rate,
		CurrentTokens:      b.currentTokens,
		LastTokenTimestamp: b.lastTokenTimestamp,
		Duration:           b.duration,
	}

	brds := b.rds.Get(ctx, identifier).Val()
	if brds == "" {
		b.rds.Set(ctx, identifier, bucket, b.duration)
		return bucket
	}

	if err := json.Unmarshal([]byte(brds), &bucket); err != nil {
		fmt.Println("invalid response from redis")
		return bucket
	}

	return bucket
}
