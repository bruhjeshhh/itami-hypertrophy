package handler

import (
	"encoding/json"
	"fmt"
	"itami-hypertrophy/internal/cache"
	"itami-hypertrophy/internal/db"
	"net/http"
	"time"
)

func GetDashboardByDate(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.Context().Value(UserEmailKey).(string)

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

// weekly tracking
func GetWeeklyDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.Context().Value(UserEmailKey).(string)

	// Determine start of week (Monday)
	startStr := r.URL.Query().Get("start")
	var weekStart time.Time
	var err error

	if startStr == "" {
		today := time.Now()
		offset := (int(today.Weekday()) + 6) % 7
		weekStart = today.AddDate(0, 0, -offset)
	} else {
		weekStart, err = time.Parse("2006-01-02", startStr)
		if err != nil {
			http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
	}

	// ✅ Redis cache key for this user + week
	cacheKey := fmt.Sprintf("weekly:%s:%s", email, weekStart.Format("2006-01-02"))

	// 1️⃣ Try to fetch from Redis first
	cached, _ := cache.Rdb.Get(cache.Ctx, cacheKey).Result()
	if cached != "" {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cached))
		return
	}

	// No cache → compute fresh
	days := []string{}
	calories := []float64{}
	protein := []float64{}
	volume := []float64{}

	var weeklyCalories float64
	var weeklyProtein float64
	var weeklyVolume float64

	for i := 0; i < 7; i++ {
		dayStart := weekStart.AddDate(0, 0, i)
		dayEnd := dayStart.Add(24 * time.Hour)
		dateStr := dayStart.Format("2006-01-02")
		days = append(days, dateStr)

		var cal, prot float64
		err := db.DB.QueryRow(`
            SELECT COALESCE(SUM(calories),0), COALESCE(SUM(protein),0)
            FROM meals
            WHERE email = $1 AND created_at >= $2 AND created_at < $3
        `, email, dayStart, dayEnd).Scan(&cal, &prot)
		if err != nil {
			http.Error(w, "DB error (meals): "+err.Error(), http.StatusInternalServerError)
			return
		}

		rows, err := db.DB.Query(`
            SELECT sets, reps, weight FROM strength_workouts
            WHERE email = $1 AND created_at >= $2 AND created_at < $3
        `, email, dayStart, dayEnd)
		if err != nil {
			http.Error(w, "DB error (workouts): "+err.Error(), http.StatusInternalServerError)
			return
		}

		var vol float64
		for rows.Next() {
			var sets, reps int
			var weight float64
			rows.Scan(&sets, &reps, &weight)
			vol += float64(sets*reps) * weight
		}
		rows.Close()

		calories = append(calories, cal)
		protein = append(protein, prot)
		volume = append(volume, vol)

		weeklyCalories += cal
		weeklyProtein += prot
		weeklyVolume += vol
	}

	// ✅ Fetch user goals
	var dailyCaloriesGoal, dailyProteinGoal, weeklyVolumeGoal float64
	err = db.DB.QueryRow(`
        SELECT daily_calories_target, daily_protein_target, weekly_volume_target
        FROM goals WHERE email = $1
    `, email).Scan(&dailyCaloriesGoal, &dailyProteinGoal, &weeklyVolumeGoal)
	if err != nil {
		dailyCaloriesGoal, dailyProteinGoal, weeklyVolumeGoal = 0, 0, 0
	}

	// Compute weekly goals
	weeklyCaloriesGoal := dailyCaloriesGoal * 7
	weeklyProteinGoal := dailyProteinGoal * 7

	// Compute % progress
	progressCalories := 0.0
	if weeklyCaloriesGoal > 0 {
		progressCalories = (weeklyCalories / weeklyCaloriesGoal) * 100
	}

	progressProtein := 0.0
	if weeklyProteinGoal > 0 {
		progressProtein = (weeklyProtein / weeklyProteinGoal) * 100
	}

	progressVolume := 0.0
	if weeklyVolumeGoal > 0 {
		progressVolume = (weeklyVolume / weeklyVolumeGoal) * 100
	}

	// ✅ Final JSON response
	response := map[string]interface{}{
		"days":     days,
		"calories": calories,
		"protein":  protein,
		"volume":   volume,
		"weekly_totals": map[string]float64{
			"calories": weeklyCalories,
			"protein":  weeklyProtein,
			"volume":   weeklyVolume,
		},
		"goals": map[string]float64{
			"weekly_calories": weeklyCaloriesGoal,
			"weekly_protein":  weeklyProteinGoal,
			"weekly_volume":   weeklyVolumeGoal,
		},
		"progress_percent": map[string]float64{
			"calories": progressCalories,
			"protein":  progressProtein,
			"volume":   progressVolume,
		},
	}

	// ✅ Cache this response in Redis for 5 min
	jsonBytes, _ := json.Marshal(response)
	cache.Rdb.Set(cache.Ctx, cacheKey, jsonBytes, 5*time.Minute)

	// ✅ Send to client
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)

}
