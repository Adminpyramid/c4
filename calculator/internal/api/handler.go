package api

import (
	"calculator/internal/db"
	"calculator/internal/models"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// SaveHandler handles POST /api/save
func SaveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}

	var calc models.Calculation
	if err := json.NewDecoder(r.Body).Decode(&calc); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if calc.SessionID == "" {
		http.Error(w, "session_id required", http.StatusBadRequest)
		return
	}

	if calc.ID == "" {
		calc.ID = uuid.New().String()
	}
	calc.CreatedAt = time.Now()

	_, err := db.DB.Exec(
		`INSERT INTO calculation_history (id, session_id, user_id, expression, result, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		calc.ID, calc.SessionID, calc.UserID, calc.Expression, calc.Result, calc.CreatedAt,
	)
	if err != nil {
		http.Error(w, "DB insert failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// HistoryHandler handles GET /api/history?session_id=...
func HistoryHandler(w http.ResponseWriter, r *http.Request) {
	session := r.URL.Query().Get("session_id")
	if session == "" {
		http.Error(w, "session_id required", http.StatusBadRequest)
		return
	}

	rows, err := db.DB.Query(`
		SELECT id, session_id, user_id, expression, result, created_at
		FROM calculation_history
		WHERE session_id = $1
		ORDER BY created_at DESC
		LIMIT 20`, session)
	if err != nil {
		http.Error(w, "DB query failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var history []models.Calculation
	for rows.Next() {
		var c models.Calculation
		if err := rows.Scan(&c.ID, &c.SessionID, &c.UserID, &c.Expression, &c.Result, &c.CreatedAt); err != nil {
			http.Error(w, "Row scan failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		history = append(history, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}
