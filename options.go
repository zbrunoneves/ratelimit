package ratelimit

import (
	"net"
	"net/http"
)

type Option func(rl *RateLimiter)

func WithKeyFunc(f func(r *http.Request) string) Option {
	return func(rl *RateLimiter) {
		rl.keyFunc = f
	}
}

func defaultKeyFunc(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}

	return r.RequestURI + host
}
