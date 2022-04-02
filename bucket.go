package ratelimit

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/go-redis/redis"
)

type Config struct {
	MaxTokens int64
	Rate      int64
	Duration  time.Duration
	Redis     *redis.Client
}

type Bucket struct {
	maxTokens          int64     `json:"max_tokens"`
	rate               int64     `json:"rate"`
	currentTokens      int64     `json:"current_tokens"`
	lastTokenTimestamp time.Time `json:"last_token_timestamp"`
	duration           time.Duration
	rds                *redis.Client
}

func New(ctx context.Context, conf Config) Bucket {
	return Bucket{
		rate:               conf.Rate,
		currentTokens:      0,
		maxTokens:          conf.MaxTokens,
		lastTokenTimestamp: time.Now(),
		rds:                conf.Redis,
		duration:           conf.Duration,
	}
}

func (b *Bucket) GetBucket(ctx context.Context, identifier string) Bucket {
	bucket := b.rds.Get(ctx, identifier).Val()
	if bucket == "" {
		b.rds.Set(ctx, identifier, *b, b.duration)
		return *b
	}
	if err := json.Unmarshal([]byte(bucket), &b); err != nil {
		fmt.Println("invalid response from redis")
		return New(ctx, Config{})
	}
	return *b
}

func (b *Bucket) IsRequestAllowed(tokens int64) bool {
	now := time.Now()
	end := time.Since(b.lastTokenTimestamp)
	tokensTobeAdded := (end.Nanoseconds() * b.rate) / 1000000000
	b.currentTokens = int64(math.Min(float64(b.currentTokens+tokensTobeAdded), float64(b.maxTokens)))
	b.lastTokenTimestamp = now

	if b.currentTokens >= tokens {
		b.currentTokens = b.currentTokens - tokens
		return true
	}
	return false
}
