package handler

import (
	"encoding/json"
	"itami-hypertrophy/internal/db"
	"net/http"
	"time"
)

type Goal struct {
	DailyCalories float64 `json:"daily_calories"`
	DailyProtein  float64 `json:"daily_protein"`
	WeeklyVolume  float64 `json:"weekly_volume"`
	UpdatedAt     string  `json:"updated_at,omitempty"`
}

// GET /goals
func GetGoals(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value(UserEmailKey()).(string)

	var g Goal
	var updatedAt time.Time
	err := db.DB.QueryRow(`
        SELECT daily_calories_target, daily_protein_target, weekly_volume_target, updated_at
        FROM goals WHERE email = $1
    `, email).Scan(&g.DailyCalories, &g.DailyProtein, &g.WeeklyVolume, &updatedAt)

	if err != nil {
		// If no record found, return defaults
		g = Goal{DailyCalories: 0, DailyProtein: 0, WeeklyVolume: 0}
	} else {
		g.UpdatedAt = updatedAt.Format(time.RFC3339)
	}

	json.NewEncoder(w).Encode(g)
}

// POST /goals
func SetGoals(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.Context().Value(UserEmailKey()).(string)

	var req Goal
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	_, err := db.DB.Exec(`
        INSERT INTO goals (email, daily_calories_target, daily_protein_target, weekly_volume_target, updated_at)
        VALUES ($1, $2, $3, $4, NOW())
        ON CONFLICT (email) DO UPDATE
        SET daily_calories_target = $2,
            daily_protein_target = $3,
            weekly_volume_target = $4,
            updated_at = NOW()
    `, email, req.DailyCalories, req.DailyProtein, req.WeeklyVolume)

	if err != nil {
		http.Error(w, "Failed to save goals: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Goals updated successfully"})
}
