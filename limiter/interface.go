package limiter

type RateLimiter interface {
	Allow(userID string) bool
}
