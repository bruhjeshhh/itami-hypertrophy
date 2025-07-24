package model

import (
	"encoding/json"
	"itami-hypertrophy/internal/db"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type creds struct {
	ID       int
	Email    string
	Password string // hash krenge isko
}

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
