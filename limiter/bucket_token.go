package limiter

import (
	"sync"
	"time"
)

type TockenBucket struct {
	Capacity     int
	RefillRate   int
	RefillPeriod time.Duration
	UserBuckets  map[string]*Bucket
	mu           sync.Mutex
}

type Bucket struct {
	Tokens         int
	LastRefillTime time.Time
	mu             sync.Mutex
}

func NewTokenBucket(capacity int, refillRate int, refillPeriod time.Duration, userBuckets map[string]*Bucket) *TockenBucket {
	return &TockenBucket{
		Capacity:     capacity,
		RefillRate:   refillRate,
		RefillPeriod: refillPeriod,
		UserBuckets:  make(map[string]*Bucket),
	}
}

func (tockenBucket *TockenBucket) getBucket(userID string) *Bucket {
	tockenBucket.mu.Lock()
	defer tockenBucket.mu.Unlock()

	// checking if bucket of user is exist then return
	if userBucket, ok := tockenBucket.UserBuckets[userID]; ok {
		return userBucket
	}

	// Init new bucket for user and return the user bucket have been init
	tockenBucket.UserBuckets[userID] = &Bucket{
		Tokens:         tockenBucket.Capacity,
		LastRefillTime: time.Now(),
	}
	return tockenBucket.UserBuckets[userID]
}

func (tockenBucket *TockenBucket) Allow(userID string) bool {
	bucket := tockenBucket.getBucket(userID)

	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	// Caculate tocken needed to refill
	currentTime := time.Now()
	elapsedTime := currentTime.Sub(bucket.LastRefillTime)
	newTocken := int(elapsedTime.Seconds()) * tockenBucket.RefillRate

	// Refill
	if newTocken > 0 {
		bucket.Tokens = min(tockenBucket.Capacity, bucket.Tokens+newTocken)
	}

	if bucket.Tokens > 0 {
		bucket.Tokens--
		return true
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
