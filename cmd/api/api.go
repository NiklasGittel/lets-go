package main

import (
	"database/sql"
	"log"
	"net/http"

	"letsgo/internal/handler"
	"letsgo/internal/repository"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "local_sales.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	saleRepo := repository.NewSQLiteRepo(db)
	h := handler.NewHandler(saleRepo)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /sales", h.GetSalesHandler)

	log.Println("Server starting on http://localhost:8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
