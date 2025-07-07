package handler

import (
	"encoding/json"
	"itami-hypertrophy/internal/db"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type creds struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func UserEmailKey() interface{} {
	return userEmailKey
}

// always remember boht maa chudi thi idhar
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var c creds
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil || c.Email == "" || c.Password == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	_, err = db.DB.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", c.Email, string(hashedPassword))
	if err != nil {
		http.Error(w, "Email already in use or DB error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}

func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Please use POST", http.StatusMethodNotAllowed)
		return
	}

	var c creds
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil || c.Email == "" || c.Password == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var storedHashedPassword string
	err = db.DB.QueryRow("SELECT password FROM users WHERE email = $1", c.Email).Scan(&storedHashedPassword)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(c.Password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// jwt —okokmaybe
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": c.Email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		http.Error(w, "Server error: missing JWT secret", http.StatusInternalServerError)
		return
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		http.Error(w, "Could not sign token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})

}
