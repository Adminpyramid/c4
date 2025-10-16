package main

import (
	"log"
	"net/http"

	"github.com/Adminpyramid/c4/calculator/internal/db"

	"github.com/Adminpyramid/c4/calculator/internal/api"
)

func main() {
	err := db.Connect("postgres", "1234", "127.0.0.1", "c4", 5432)
	if err != nil {
		log.Fatal("[DB] Connection failed:", err)
	}

	http.HandleFunc("/api/save", api.SaveHandler)
	http.HandleFunc("/api/history", api.HistoryHandler)

	log.Println("[Server] Running on http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
