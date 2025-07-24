package handler

import (
	"encoding/json"
	"itami-hypertrophy/internal/db"
	"net/http"
)

type Goals struct {
	DailyCalories       int     `json:"daily_calories"`
	DailyProtein        float64 `json:"daily_protein"`
	WeeklyWorkoutVolume int     `json:"weekly_workout_volume"`
}

// GET /goals → fetch current goals
func GetGoals(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value(UserEmailKey).(string)

	var g Goals
	err := db.DB.QueryRow(`
        SELECT daily_calories, daily_protein, weekly_workout_volume 
        FROM goals WHERE email = $1
    `, email).Scan(&g.DailyCalories, &g.DailyProtein, &g.WeeklyWorkoutVolume)

	if err != nil {
		http.Error(w, "No goals found. Please set them first.", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(g)
}

// POST /goals/set → update or create goals
func SetGoals(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value(UserEmailKey).(string)

	var g Goals
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	_, err := db.DB.Exec(`
        INSERT INTO goals (email, daily_calories, daily_protein, weekly_workout_volume)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (email) DO UPDATE SET
        daily_calories = EXCLUDED.daily_calories,
        daily_protein = EXCLUDED.daily_protein,
        weekly_workout_volume = EXCLUDED.weekly_workout_volume
    `, email, g.DailyCalories, g.DailyProtein, g.WeeklyWorkoutVolume)

	if err != nil {
		http.Error(w, "Failed to save goals: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Goals updated successfully"))
}
