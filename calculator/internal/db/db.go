package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// Connect initializes the PostgreSQL connection pool.
func Connect(user, password, host, dbname string, port int) error {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		user, password, host, port, dbname,
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open DB: %w", err)
	}

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	log.Println("[DB] Connected successfully")
	return nil
}
