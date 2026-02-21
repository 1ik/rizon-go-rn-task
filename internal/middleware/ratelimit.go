package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"rizon-test-task/internal/config"
)

// windowEntry holds the count and window end for a single key (e.g. IP).
type windowEntry struct {
	count    int
	windowEnd time.Time
}

// ipLimiter is a fixed-window rate limiter keyed by IP.
type ipLimiter struct {
	mu       sync.Mutex
	entries  map[string]*windowEntry
	limit    int
	window   time.Duration
}

func newIPLimiter(cfg *config.RateLimitConfig) *ipLimiter {
	return &ipLimiter{
		entries: make(map[string]*windowEntry),
		limit:   cfg.RequestsPerWindow,
		window:  cfg.Window,
	}
}

// allow returns true if the key is within limit, false if rate limited.
func (l *ipLimiter) allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	e, ok := l.entries[key]
	if !ok || now.After(e.windowEnd) {
		l.entries[key] = &windowEntry{count: 1, windowEnd: now.Add(l.window)}
		return true
	}
	if e.count >= l.limit {
		return false
	}
	e.count++
	return true
}

// clientIP returns the client IP from the request (X-Forwarded-For or RemoteAddr).
func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// First value is the client IP when behind a single proxy
		if i := strings.Index(xff, ","); i >= 0 {
			return strings.TrimSpace(xff[:i])
		}
		return strings.TrimSpace(xff)
	}
	// RemoteAddr is "ip:port"
	addr := r.RemoteAddr
	if i := strings.LastIndex(addr, ":"); i >= 0 {
		return addr[:i]
	}
	return addr
}

// RateLimit returns a middleware that limits requests per client IP.
// OPTIONS requests are not counted (CORS preflight). Returns 429 with Retry-After when exceeded.
func RateLimit(cfg *config.RateLimitConfig) func(http.Handler) http.Handler {
	limiter := newIPLimiter(cfg)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}
			ip := clientIP(r)
			if !limiter.allow(ip) {
				w.Header().Set("Retry-After", strconv.Itoa(int(cfg.Window.Seconds())))
				http.Error(w, "rate limit exceeded, try again later", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
