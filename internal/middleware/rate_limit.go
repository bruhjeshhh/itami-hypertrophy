package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"itami-hypertrophy/internal/cache"
)

func RateLimitMiddleware(limit int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Use user email (from JWT) if available, else IP
			userKey := extractUserKey(r)
			if userKey == "" {
				userKey = r.RemoteAddr // fallback to IP
			}

			redisKey := fmt.Sprintf("ratelimit:%s:%s", r.URL.Path, userKey)

			// Increment request count
			count, err := cache.Rdb.Incr(cache.Ctx, redisKey).Result()
			if err != nil {
				http.Error(w, "Redis error", http.StatusInternalServerError)
				return
			}

			// Set TTL only on first request
			if count == 1 {
				cache.Rdb.Expire(cache.Ctx, redisKey, window)
			}

			// Check limit
			if count > int64(limit) {
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(" Too many requests, slow down!"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Extract user email from Authorization header if present
func extractUserKey(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if auth != "" && strings.HasPrefix(auth, "Bearer ") {
		// You could decode JWT here and get email claim
		return auth[7:] // For now just use token as key
	}
	return ""
}
