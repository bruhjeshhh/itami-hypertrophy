package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/cors"

	"itami-hypertrophy/internal/cache"
	"itami-hypertrophy/internal/db"
	"itami-hypertrophy/internal/handler"
)

func main() {
	godotenv.Load()
	fmt.Println("Using DB_URL:", os.Getenv("DB_URL"))
	db.Connect()
	cache.InitRedis()

	mux := http.NewServeMux()

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong")
	})
	mux.HandleFunc("/register", handler.Register)
	mux.HandleFunc("/login", handler.Login)
	mux.HandleFunc("/profile", handler.JWTMiddleware(func(w http.ResponseWriter, r *http.Request) {
		email := r.Context().Value(handler.UserEmailKey).(string)
		w.Write([]byte("Hello, " + email + "! This is your profile."))
	}))
	mux.HandleFunc("/log-calories", handler.JWTMiddleware(handler.LogCalories))
	mux.HandleFunc("/meals", handler.JWTMiddleware(handler.GetMeals))
	mux.HandleFunc("/meals/today", handler.JWTMiddleware(handler.GetTodayMeals))
	mux.HandleFunc("/log-strength", handler.JWTMiddleware(handler.LogStrengthWorkout))
	mux.HandleFunc("/dashboard", handler.JWTMiddleware(handler.GetDashboardByDate))
	mux.HandleFunc("/dashboard/weekly", handler.JWTMiddleware(handler.GetWeeklyDashboard))
	mux.HandleFunc("/goals", handler.JWTMiddleware(handler.GetGoals))
	mux.HandleFunc("/goals/set", handler.JWTMiddleware(handler.SetGoals))

	// ✅ Enable CORS for frontend
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	fmt.Println("Running backend on :8080")
	http.ListenAndServe(":8080", c.Handler(mux))
}
