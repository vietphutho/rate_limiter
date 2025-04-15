package limiter

import (
	"sync"
	"time"
)

type FixedWindow struct {
	Limit       int
	WindowSize  time.Duration
	UserWindows map[string]*Window
	mu          sync.Mutex
}

type Window struct {
	StartTime time.Time
	Count     int
}

func NewFixedWindow(limit int, windowSize time.Duration) *FixedWindow {
	return &FixedWindow{
		Limit:       limit,
		WindowSize:  windowSize,
		UserWindows: make(map[string]*Window),
	}
}

func (fixedWindow *FixedWindow) Allow(userID string) bool {
	fixedWindow.mu.Lock()
	defer fixedWindow.mu.Unlock()

	currentTime := time.Now()
	window, exists := fixedWindow.UserWindows[userID]

	if !exists || currentTime.Sub(window.StartTime) >= fixedWindow.WindowSize {
		fixedWindow.UserWindows[userID] = &Window{
			StartTime: currentTime,
			Count:     1,
		}
		return true
	}

	if window.Count < fixedWindow.Limit {
		window.Count++
		return true
	}
	return false
}
