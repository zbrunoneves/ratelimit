package ratelimit

import (
	"sync"
	"time"
)

type slidingWindow struct {
	limit  int
	window time.Duration
	store  map[string][]time.Time
	mu     *sync.Mutex
}

func newSlidingWindow(limit int, window time.Duration) *slidingWindow {
	return &slidingWindow{
		limit:  limit,
		window: window,
		store:  map[string][]time.Time{},
		mu:     &sync.Mutex{},
	}
}

func (sw *slidingWindow) allow(key string) (bool, time.Duration) {
	now := time.Now()

	sw.mu.Lock()
	defer sw.mu.Unlock()

	_, ok := sw.store[key]
	if !ok {
		sw.store[key] = []time.Time{now}
		return true, 0
	}

	start := now.Add(-sw.window)
	entry := sw.store[key]

	var i int
	for i < len(entry) {
		if entry[i].After(start) {
			break
		}
		i++
	}

	fresh := entry[i:]
	if len(fresh) >= sw.limit {
		return false, fresh[0].Sub(start)
	}

	sw.store[key] = append(fresh, now)

	return true, 0
}
