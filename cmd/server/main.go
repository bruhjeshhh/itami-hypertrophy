package main

import (
	"fmt"
	"net/http"
	"os"

	"itami-hypertrophy/internal/cache"
	"itami-hypertrophy/internal/db"
	"itami-hypertrophy/internal/handler"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	fmt.Println("Using DB_URL:", os.Getenv("DB_URL"))

	// ✅ Call your existing DB connector
	db.Connect()

	// ✅ Init Redis before routes
	cache.InitRedis()
	fmt.Println("Connected to Redis")

	// ✅ Public routes
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong")
	})
	http.HandleFunc("/register", handler.Register)
	http.HandleFunc("/login", handler.Login)

	// ✅ Authenticated routes
	http.HandleFunc("/profile", handler.JWTMiddleware(func(w http.ResponseWriter, r *http.Request) {
		email := r.Context().Value(handler.UserEmailKey).(string)
		w.Write([]byte("Hello, " + email + "! This is your profile."))
	}))

	http.HandleFunc("/log-calories", handler.JWTMiddleware(handler.LogCalories))
	http.HandleFunc("/meals", handler.JWTMiddleware(handler.GetMeals))
	http.HandleFunc("/meals/today", handler.JWTMiddleware(handler.GetTodayMeals))
	http.HandleFunc("/log-strength", handler.JWTMiddleware(handler.LogStrengthWorkout))
	http.HandleFunc("/dashboard", handler.JWTMiddleware(handler.GetDashboardByDate))
	http.HandleFunc("/dashboard/weekly", handler.JWTMiddleware(handler.GetWeeklyDashboard))
	http.HandleFunc("/goals", handler.JWTMiddleware(handler.GetGoals))
	http.HandleFunc("/goals/set", handler.JWTMiddleware(handler.SetGoals))

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
