package ratelimit

import (
	"errors"
	"net/http"
	"time"
)

const (
	AlgorithmFixedWindow   = "fixed-window"
	AlgorithmSlidingWindow = "sliding-window"
)

type limiter interface {
	allow(key string) (bool, time.Duration)
}

type RateLimiter struct {
	keyFunc func(r *http.Request) string
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
		return nil, errors.New("rate limit: unknown algorithm " + algorithm)
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
	key := rl.keyFunc(r)

	return rl.limiter.allow(key)
}
