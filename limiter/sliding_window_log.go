package limiter

import (
	"sync"
	"time"
)

type SlidingWindowLog struct {
	Limit      int
	WindowSize time.Duration
	UserLogs   map[string][]time.Time
	mu         sync.Mutex
}

func NewSlidingWindownLog(limit int, windowSize time.Duration) *SlidingWindowLog {
	return &SlidingWindowLog{
		Limit:      limit,
		WindowSize: windowSize,
		UserLogs:   make(map[string][]time.Time),
	}
}

func (slidingWindowLog *SlidingWindowLog) Allow(userID string) bool {
	slidingWindowLog.mu.Lock()
	defer slidingWindowLog.mu.Unlock()

	currentTime := time.Now()
	userLogs := slidingWindowLog.UserLogs[userID]
	filterLogs := []time.Time{}

	for _, t := range userLogs {
		if currentTime.Sub(t) < slidingWindowLog.WindowSize {
			filterLogs = append(filterLogs, t)
		}
	}

	if len(filterLogs) < slidingWindowLog.Limit {
		filterLogs = append(filterLogs, currentTime)
		slidingWindowLog.UserLogs[userID] = filterLogs
		return true
	}

	return false
}
