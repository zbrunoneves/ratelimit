package ratelimit

import (
	"sync"
	"time"
)

type slidingWindow struct {
	limit  int
	window time.Duration
	store  map[string]fixedWindowEntry
	mu     *sync.Mutex
}

func newSlidingWindow(limit int, window time.Duration) *slidingWindow {
	return &slidingWindow{
		limit:  limit,
		window: window,
		store:  map[string]fixedWindowEntry{},
		mu:     &sync.Mutex{},
	}
}

func (sw *slidingWindow) allow(key string) (bool, time.Duration) {
	return false, maxDuration
}
