package config

import (
	"strconv"
	"time"
)

// RateLimitConfig holds rate limiting configuration.
type RateLimitConfig struct {
	// RequestsPerWindow is the maximum number of requests allowed per IP per window.
	RequestsPerWindow int
	// Window is the duration of each rate limit window (e.g. 1 minute).
	Window time.Duration
}

// GetRateLimitConfig returns rate limit configuration from environment or defaults.
// Env: RATE_LIMIT_REQUESTS (default 60), RATE_LIMIT_WINDOW_SECONDS (default 60).
func GetRateLimitConfig() *RateLimitConfig {
	requests := 60
	if n, err := strconv.Atoi(getEnv("RATE_LIMIT_REQUESTS", "60")); err == nil && n > 0 {
		requests = n
	}
	windowSec := 60
	if n, err := strconv.Atoi(getEnv("RATE_LIMIT_WINDOW_SECONDS", "60")); err == nil && n > 0 {
		windowSec = n
	}
	return &RateLimitConfig{
		RequestsPerWindow: requests,
		Window:            time.Duration(windowSec) * time.Second,
	}
}
