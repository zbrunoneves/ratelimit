package ratelimit

import (
	"errors"
	"net/http"
	"time"
)

const (
	AlgorithmFixedWindow   = "fixed-window"
	AlgorithmSlidingWindow = "sliding-window"

	maxDuration = 1<<63 - 1
)

type limiter interface {
	allow(key string) (bool, time.Duration)
}

type RateLimiter struct {
	keyFunc func(r *http.Request) (string, error)
	limiter limiter
}

func New(limit int, window time.Duration, algorithm string, opts ...Option) (*RateLimiter, error) {
	var l limiter

	switch algorithm {
	case AlgorithmFixedWindow:
		l = newFixedWindow(limit, window)
	case AlgorithmSlidingWindow:
		l = newSlidingWindow(limit, window)
	default:
		return nil, errors.New("unknown algorithm: " + algorithm)
	}

	rl := &RateLimiter{
		keyFunc: defaultKeyFunc,
		limiter: l,
	}

	for _, opt := range opts {
		opt(rl)
	}

	return rl, nil
}

func (rl RateLimiter) Allow(r *http.Request) (bool, time.Duration) {
	key, err := rl.keyFunc(r)
	if err != nil {
		// should log error
		return false, maxDuration
	}

	return rl.limiter.allow(key)
}
