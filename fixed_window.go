package ratelimit

import (
	"sync"
	"time"
)

type fixedWindow struct {
	limit  int
	window time.Duration
	store  map[string]fixedWindowEntry
	mu     *sync.Mutex
}

type fixedWindowEntry struct {
	count     int
	windowEnd time.Time
}

func newFixedWindow(limit int, window time.Duration) *fixedWindow {
	return &fixedWindow{
		limit:  limit,
		window: window,
		store:  map[string]fixedWindowEntry{},
		mu:     &sync.Mutex{},
	}
}

func (fw *fixedWindow) allow(key string) (bool, time.Duration) {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	now := time.Now()
	entry := fw.store[key]

	if now.After(entry.windowEnd) {
		fw.store[key] = fixedWindowEntry{
			count:     1,
			windowEnd: now.Add(fw.window),
		}

		return true, 0
	}

	if entry.count >= fw.limit {
		return false, entry.windowEnd.Sub(now)
	}

	entry.count++
	fw.store[key] = entry

	return true, 0
}
