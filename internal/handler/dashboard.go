package handler

import (
	"encoding/json"
	"itami-hypertrophy/internal/db"
	"net/http"
	"time"
)

func GetDashboardByDate(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.Context().Value(UserEmailKey()).(string)

	dateStr := r.URL.Query().Get("date")
	var targetDate time.Time
	var err error

	if dateStr == "" {
		targetDate = time.Now()
	} else {
		targetDate, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
	}

	start := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, targetDate.Location())
	end := start.Add(24 * time.Hour)

	// gets meals
	mealsRows, err := db.DB.Query(`
		SELECT description, calories, protein, carbs, fat, created_at
		FROM meals
		WHERE email = $1 AND created_at >= $2 AND created_at < $3
		ORDER BY created_at ASC
	`, email, start, end)
	if err != nil {
		http.Error(w, "DB error (meals): "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer mealsRows.Close()

	type Meal struct {
		Description string  `json:"description"`
		Calories    float64 `json:"calories"`
		Protein     float64 `json:"protein"`
		Carbs       float64 `json:"carbs"`
		Fat         float64 `json:"fat"`
		LoggedAt    string  `json:"logged_at"`
	}

	var meals []Meal
	var totalCalories, totalProtein, totalCarbs, totalFat float64

	for mealsRows.Next() {
		var m Meal
		var createdAt time.Time
		err := mealsRows.Scan(&m.Description, &m.Calories, &m.Protein, &m.Carbs, &m.Fat, &createdAt)
		if err != nil {
			http.Error(w, "Row scan failed (meals)", http.StatusInternalServerError)
			return
		}
		m.LoggedAt = createdAt.Format(time.RFC3339)
		meals = append(meals, m)

		totalCalories += m.Calories
		totalProtein += m.Protein
		totalCarbs += m.Carbs
		totalFat += m.Fat
	}

	// gets workouts

	workoutRows, err := db.DB.Query(`
		SELECT exercise, sets, reps, weight, created_at
		FROM strength_workouts
		WHERE email = $1 AND created_at >= $2 AND created_at < $3
		ORDER BY created_at ASC
	`, email, start, end)
	if err != nil {
		http.Error(w, "DB error (workouts): "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer workoutRows.Close()

	type Workout struct {
		Exercise string  `json:"exercise"`
		Sets     int     `json:"sets"`
		Reps     int     `json:"reps"`
		Weight   float64 `json:"weight"`
		LoggedAt string  `json:"logged_at"`
		Volume   float64 `json:"volume"`
	}

	var workouts []Workout
	var totalSets, totalReps int
	var totalVolume float64

	for workoutRows.Next() {
		var workout Workout
		var createdAt time.Time
		err := workoutRows.Scan(&workout.Exercise, &workout.Sets, &workout.Reps, &workout.Weight, &createdAt)
		if err != nil {
			http.Error(w, "Row scan failed (workouts)", http.StatusInternalServerError)
			return
		}
		workout.LoggedAt = createdAt.Format(time.RFC3339)
		workout.Volume = float64(workout.Sets*workout.Reps) * workout.Weight
		workouts = append(workouts, workout)

		totalSets += workout.Sets
		totalReps += workout.Sets * workout.Reps
		totalVolume += workout.Volume
	}

	summary := map[string]interface{}{
		"calories":     totalCalories,
		"protein":      totalProtein,
		"carbs":        totalCarbs,
		"fat":          totalFat,
		"total_sets":   totalSets,
		"total_reps":   totalReps,
		"total_volume": totalVolume,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"meals":    meals,
		"workouts": workouts,
		"summary":  summary,
	})
}
