package handler

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

// ✅ Only ONE declaration
var UserEmailKey contextKey = "userEmail"

func JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			http.Error(w, "Server misconfigured: missing JWT secret", http.StatusInternalServerError)
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		email, _ := claims["email"].(string)
		if email == "" {
			http.Error(w, "Missing email in token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserEmailKey, email)
		next(w, r.WithContext(ctx))
	}
}

func GetUserEmail(r *http.Request) string {
	if email, ok := r.Context().Value(UserEmailKey).(string); ok {
		return email
	}
	return ""
}
