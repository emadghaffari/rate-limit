package ratelimit

import (
	"context"
	"encoding/json"
	"math"
	"time"

	"github.com/go-redis/redis/v8"
)

type Bucket struct {
	MaxTokens          int64         `json:"max_tokens"`
	CurrentTokens      int64         `json:"current_tokens"`
	LastTokenTimestamp time.Time     `json:"last_token_timestamp"`
	rds                *redis.Client `json:"-"`
	Identifier         string        `json:"identifier"`
	Duration           time.Duration `json:"duration"`
}

// IsRequestAllowed: check the request with the tokens
func (b *Bucket) IsRequestAllowed(ctx context.Context, tokens int64) bool {
	tokensTobeAdded := time.Since(b.LastTokenTimestamp).Nanoseconds() / 10e9
	b.CurrentTokens = int64(math.Min(float64(b.CurrentTokens+tokensTobeAdded), float64(b.MaxTokens)))
	b.LastTokenTimestamp = time.Now()

	if b.CurrentTokens < b.MaxTokens {
		b.IncreaseToken(ctx, tokens)
		return true
	}

	return false
}

// IncreaseToken: increase current token number in redis
func (b *Bucket) IncreaseToken(ctx context.Context, tokens int64) error {
	b.CurrentTokens = b.CurrentTokens + tokens
	if b.CurrentTokens > b.MaxTokens {
		b.CurrentTokens = b.MaxTokens
	}

	bts, _ := json.Marshal(b)
	return b.rds.Set(context.Background(), b.Identifier, bts, b.Duration).Err()
}

// DecreaseToken: decrease current token number in redis
func (b *Bucket) DecreaseToken(ctx context.Context, tokens int64) error {
	b.CurrentTokens = b.CurrentTokens - tokens
	if b.CurrentTokens < 0 {
		b.CurrentTokens = 0
	}

	bts, _ := json.Marshal(b)
	return b.rds.Set(context.Background(), b.Identifier, bts, b.Duration).Err()
}
