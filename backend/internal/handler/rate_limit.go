package handler

import (
	"fmt"
	"net/http"
	"time"

	"itami-hypertrophy/internal/cache"
)

func RateLimit(limit int, window time.Duration, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// ✅ fetch email from JWT context
		email, ok := r.Context().Value(UserEmailKey).(string)
		if !ok || email == "" {
			http.Error(w, "Unauthorized (no email in context)", http.StatusUnauthorized)
			return
		}

		route := r.URL.Path
		key := fmt.Sprintf("rate:%s:%s", email, route)

		// Redis increment
		count, err := cache.Rdb.Incr(cache.Ctx, key).Result()
		if err != nil {
			http.Error(w, "Rate limiter error", http.StatusInternalServerError)
			return
		}

		if count == 1 {
			cache.Rdb.Expire(cache.Ctx, key, window)
		}

		if int(count) > limit {
			ttl, _ := cache.Rdb.TTL(cache.Ctx, key).Result()
			http.Error(w, fmt.Sprintf("Rate limit exceeded. Try again in %.0f seconds", ttl.Seconds()), http.StatusTooManyRequests)
			return
		}

		next(w, r)
	}
}
