// package middleware

// import (
// 	"fmt"
// 	"net/http"
// 	"strings"
// 	"time"

// 	"itami-hypertrophy/internal/cache"
// )

// func RateLimitMiddleware(limit int, window time.Duration) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 			userKey := extractUserKey(r)
// 			if userKey == "" {
// 				userKey = r.RemoteAddr
// 			}

// 			redisKey := fmt.Sprintf("ratelimit:%s:%s", r.URL.Path, userKey)

// 			count, err := cache.Rdb.Incr(cache.Ctx, redisKey).Result()
// 			if err != nil {
// 				http.Error(w, "Redis error", http.StatusInternalServerError)
// 				return
// 			}

// 			if count == 1 {
// 				cache.Rdb.Expire(cache.Ctx, redisKey, window)
// 			}

// 		t
// 			if count > int64(limit) {
// 				w.WriteHeader(http.StatusTooManyRequests)
// 				w.Write([]byte(" Too many requests, slow down!"))
// 				return
// 			}

// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }

// func extractUserKey(r *http.Request) string {
// 	auth := r.Header.Get("Authorization")
// 	if auth != "" && strings.HasPrefix(auth, "Bearer ") {

// 		return auth[7:]
// 	}
// 	return ""
// }
