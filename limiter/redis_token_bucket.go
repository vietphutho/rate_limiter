package limiter

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisTokenBucket struct {
	RedisClient *redis.Client
	Capacity    int
	RefillRate  int // tokens per second
	KeyPrefix   string
	mu          sync.Mutex
}

func NewRedisTokenBucket(redisClient *redis.Client, capacity, refillRate int, keyPrefix string) *RedisTokenBucket {
	return &RedisTokenBucket{
		RedisClient: redisClient,
		Capacity:    capacity,
		RefillRate:  refillRate,
		KeyPrefix:   keyPrefix,
	}
}

func (r *RedisTokenBucket) Allow(userID string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	ctx := context.Background()
	bucketKey := fmt.Sprintf("%s:bucket:%s", r.KeyPrefix, userID)
	lastRefillKey := fmt.Sprintf("%s:lastRefill:%s", r.KeyPrefix, userID)

	currentTime := time.Now().Unix()
	lastRefillStr, err := r.RedisClient.Get(ctx, lastRefillKey).Result()
	lastRefill := currentTime
	if err == nil {
		if t, err := strconv.ParseInt(lastRefillStr, 10, 64); err == nil {
			lastRefill = t
		}
	}

	tokensStr, err := r.RedisClient.Get(ctx, bucketKey).Result()
	var tokens int
	if err == redis.Nil {
		tokens = r.Capacity
	} else if err == nil {
		tokens, _ = strconv.Atoi(tokensStr)
	} else {
		return false
	}

	elapsed := currentTime - lastRefill
	tokens += int(elapsed) * r.RefillRate
	if tokens > r.Capacity {
		tokens = r.Capacity
	}

	if tokens > 0 {
		tokens--
		r.RedisClient.Set(ctx, bucketKey, tokens, 0)
		r.RedisClient.Set(ctx, lastRefillKey, currentTime, 0)
		return true
	}

	return false
}
