package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"itami-hypertrophy/internal/cache"
	"itami-hypertrophy/internal/db"
	"itami-hypertrophy/internal/handler"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	godotenv.Load()
	fmt.Println("Using DB_URL:", os.Getenv("DB_URL"))

	// Connect to Postgres
	db.Connect()

	// Init Redis
	cache.InitRedis()

	// Health check route
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong")
	})

	// Auth routes
	http.HandleFunc("/register", handler.Register)
	http.HandleFunc("/login", handler.Login)

	// Profile (protected)
	http.HandleFunc("/profile", handler.JWTMiddleware(func(w http.ResponseWriter, r *http.Request) {
		email := r.Context().Value(handler.UserEmailKey).(string)
		w.Write([]byte("Hello, " + email + "! This is your profile."))
	}))

	// --- ✅ Protected Routes ---
	// Meals
	http.HandleFunc("/log-calories",
		handler.JWTMiddleware(
			handler.RateLimit(5, time.Minute, handler.LogCalories), // ⏳ 5 req/min
		),
	)
	http.HandleFunc("/meals", handler.JWTMiddleware(handler.GetMeals))
	http.HandleFunc("/meals/today", handler.JWTMiddleware(handler.GetTodayMeals))

	// Strength workouts
	http.HandleFunc("/log-strength",
		handler.JWTMiddleware(
			handler.RateLimit(5, time.Minute, handler.LogStrengthWorkout), // ⏳ 5 req/min
		),
	)

	// Dashboard & Goals
	http.HandleFunc("/dashboard", handler.JWTMiddleware(handler.GetDashboardByDate))
	http.HandleFunc("/dashboard/weekly", handler.JWTMiddleware(handler.GetWeeklyDashboard))
	http.HandleFunc("/goals", handler.JWTMiddleware(handler.GetGoals))
	http.HandleFunc("/goals/set",
		handler.JWTMiddleware(
			handler.RateLimit(3, time.Minute, handler.SetGoals), // ⏳ 3 req/min
		),
	)

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
