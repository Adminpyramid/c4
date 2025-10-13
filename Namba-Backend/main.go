package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Calculation struct {
	ID         string    `json:"id"`
	SessionID  string    `json:"session_id"`
	UserID     string    `json:"user_id"`
	Expression string    `json:"expression"`
	Result     string    `json:"result"`
	CreatedAt  time.Time `json:"created_at"`
}

var db *sql.DB

func main() {
	// ⚙️ Change this to your real credentials
	connStr := "postgres://postgres:yourpassword@127.0.0.1:5432/yourdb?sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("DB connection failed:", err)
	}
	log.Println("Connected to PostgreSQL")

	http.HandleFunc("/api/save", saveHandler)
	http.HandleFunc("/api/history", historyHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}

	var calc Calculation
	if err := json.NewDecoder(r.Body).Decode(&calc); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	_, err := db.Exec(
		`INSERT INTO calculation_history (session_id, user_id, expression, result)
		 VALUES ($1, $2, $3, $4)`,
		calc.SessionID, calc.UserID, calc.Expression, calc.Result,
	)
	if err != nil {
		http.Error(w, "DB insert failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func historyHandler(w http.ResponseWriter, r *http.Request) {
	session := r.URL.Query().Get("session_id")
	if session == "" {
		http.Error(w, "session_id required", http.StatusBadRequest)
		return
	}

	rows, err := db.Query(`
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

	var history []Calculation
	for rows.Next() {
		var c Calculation
		if err := rows.Scan(&c.ID, &c.SessionID, &c.UserID, &c.Expression, &c.Result, &c.CreatedAt); err != nil {
			http.Error(w, "Row scan failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		history = append(history, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}
