package main

import (
	"fmt"
	"net/http"
	"os"

	"itami-hypertrophy/internal/db"
	"itami-hypertrophy/internal/handler"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	fmt.Println("Using DB_URL:", os.Getenv("DB_URL"))
	db.Connect()

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong")
	})

	http.HandleFunc("/register", handler.Register)
	http.HandleFunc("/login", handler.Login)
	http.HandleFunc("/profile", handler.JWTMiddleware(func(w http.ResponseWriter, r *http.Request) {
		email := r.Context().Value(handler.UserEmailKey()).(string)
		w.Write([]byte("Hello, " + email + "! This is your profile."))
	}))

	fmt.Println("Server starting on port 8080...")
	http.ListenAndServe(":8080", nil)
}
