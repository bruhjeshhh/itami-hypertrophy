package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"itami-hypertrophy/internal/db"
	"net/http"
	"os"
	"time"
)

type calorieRequest struct {
	Description string `json:"description"`
}

type NutritionResult struct {
	Calories float64
	Protein  float64
	Carbs    float64
	Fat      float64
}

func fetchNutritionFromNutritionix(description string) (NutritionResult, error) {
	appID := os.Getenv("NUTRITIONIX_APP_ID")
	appKey := os.Getenv("NUTRITIONIX_APP_KEY")

	reqBody := map[string]string{"query": description}
	jsonData, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", "https://trackapi.nutritionix.com/v2/natural/nutrients", bytes.NewBuffer(jsonData))
	if err != nil {
		return NutritionResult{}, err
	}

	req.Header.Set("x-app-id", appID)
	req.Header.Set("x-app-key", appKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return NutritionResult{}, err
	}
	defer resp.Body.Close()

	var response struct {
		Foods []struct {
			Calories float64 `json:"nf_calories"`
			Protein  float64 `json:"nf_protein"`
			Carbs    float64 `json:"nf_total_carbohydrate"`
			Fat      float64 `json:"nf_total_fat"`
		} `json:"foods"`
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil || len(response.Foods) == 0 {
		return NutritionResult{}, fmt.Errorf("failed to parse response")
	}

	var total NutritionResult
	for _, food := range response.Foods {
		total.Calories += food.Calories
		total.Protein += food.Protein
		total.Carbs += food.Carbs
		total.Fat += food.Fat
	}
	return total, nil
}

func LogCalories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "post onli", http.StatusMethodNotAllowed)
		return
	}

	var req calorieRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Description == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	email := r.Context().Value(UserEmailKey()).(string)

	nutrition, err := fetchNutritionFromNutritionix(req.Description)
	if err != nil {
		http.Error(w, "ailed to fetch nutrition "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.DB.Exec(`
		INSERT INTO meals (email, description, calories, protein, carbs, fat)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, email, req.Description, nutrition.Calories, nutrition.Protein, nutrition.Carbs, nutrition.Fat)

	if err != nil {
		http.Error(w, "failed to save "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"email":       email,
		"description": req.Description,
		"calories":    nutrition.Calories,
		"protein":     nutrition.Protein,
		"carbs":       nutrition.Carbs,
		"fat":         nutrition.Fat,
	})
}

// gpt-ed coz didnt know how to use the time wala thing also sleepy
func GetMeals(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	} //ykwitmeans

	email := r.Context().Value(UserEmailKey()).(string) // jwt se user ki info nikali aur ab wahi dikhaenge jo user hai not kisi aur ka

	rows, err := db.DB.Query(`
		SELECT description, calories, protein, carbs, fat, created_at
		FROM meals
		WHERE email = $1
		ORDER BY created_at DESC
	`, email) //db se wo uthaya jo chahiye aur usko aaj ke hisab se sort kia
	if err != nil {
		http.Error(w, "Failed to fetch meals: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Meal struct {
		Description string  `json:"description"`
		Calories    float64 `json:"calories"`
		Protein     float64 `json:"protein"`
		Carbs       float64 `json:"carbs"`
		Fat         float64 `json:"fat"`
		LoggedAt    string  `json:"logged_at"`
	} //display krne ko

	var meals []Meal
	for rows.Next() {
		var m Meal
		var createdAt time.Time
		err := rows.Scan(&m.Description, &m.Calories, &m.Protein, &m.Carbs, &m.Fat, &createdAt)
		if err != nil {
			http.Error(w, "Row scan failed", http.StatusInternalServerError)
			return
		}
		m.LoggedAt = createdAt.Format(time.RFC3339)
		meals = append(meals, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(meals)
}

// again gpted-aaj ka dikhaega
func GetTodayMeals(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.Context().Value(UserEmailKey()).(string)

	// time uthaya
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.Add(24 * time.Hour)

	rows, err := db.DB.Query(`
		SELECT description, calories, protein, carbs, fat, created_at
		FROM meals
		WHERE email = $1 AND created_at >= $2 AND created_at < $3
		ORDER BY created_at ASC
	`, email, start, end)
	if err != nil {
		http.Error(w, "Failed to fetch meals: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Meal struct {
		Description string  `json:"description"`
		Calories    float64 `json:"calories"`
		Protein     float64 `json:"protein"`
		Carbs       float64 `json:"carbs"`
		Fat         float64 `json:"fat"`
		LoggedAt    string  `json:"logged_at"`
	}

	var meals []Meal
	var total NutritionResult

	for rows.Next() {
		var m Meal
		var createdAt time.Time
		err := rows.Scan(&m.Description, &m.Calories, &m.Protein, &m.Carbs, &m.Fat, &createdAt)
		if err != nil {
			http.Error(w, "Row scan failed", http.StatusInternalServerError)
			return
		}
		m.LoggedAt = createdAt.Format(time.RFC3339)
		meals = append(meals, m)

		total.Calories += m.Calories
		total.Protein += m.Protein
		total.Carbs += m.Carbs
		total.Fat += m.Fat
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"meals":   meals,
		"summary": total,
	})
}
