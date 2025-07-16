package handler

import (
	"encoding/json"
	"itami-hypertrophy/internal/db"
	"net/http"
)

type StrengthWorkoutRequest struct {
	Exercise string  `json:"exercise"`
	Sets     int     `json:"sets"`
	Reps     int     `json:"reps"`
	Weight   float64 `json:"weight"`
}

func LogStrengthWorkout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.Context().Value(UserEmailKey()).(string)

	var req StrengthWorkoutRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Exercise == "" || req.Sets <= 0 || req.Reps <= 0 || req.Weight < 0 {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec(`
		INSERT INTO strength_workouts (email, exercise, sets, reps, weight)
		VALUES ($1, $2, $3, $4, $5)
	`, email, req.Exercise, req.Sets, req.Reps, req.Weight)

	if err != nil {
		http.Error(w, "Failed to save workout: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Workout logged successfully",
	})
}
