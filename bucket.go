package ratelimit

import (
	"math"
	"time"

	"github.com/go-redis/redis"
)

type Bucket struct {
	MaxTokens          int64         `json:"max_tokens"`
	Rate               int64         `json:"rate"`
	CurrentTokens      int64         `json:"current_tokens"`
	LastTokenTimestamp time.Time     `json:"last_token_timestamp"`
	Duration           time.Duration `json:"duration"`
	rds                *redis.Client
}

func (b *Bucket) IsRequestAllowed(tokens int64) bool {
	now := time.Now()
	end := time.Since(b.LastTokenTimestamp)
	tokensTobeAdded := (end.Nanoseconds() * b.Rate) / 1000000000
	b.CurrentTokens = int64(math.Min(float64(b.CurrentTokens+tokensTobeAdded), float64(b.MaxTokens)))
	b.LastTokenTimestamp = now

	if b.CurrentTokens >= tokens {
		b.CurrentTokens = b.CurrentTokens - tokens
		return true
	}
	return false
}
