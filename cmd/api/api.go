package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"letsgo/internal/handler"
	"letsgo/internal/repository"
	"letsgo/internal/service"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "local_sales.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	saleRepo := repository.NewSQLiteRepo(db, 2*time.Second)
	saleService := service.NewSalesService(saleRepo, time.Minute)
	h := handler.NewHandler(saleService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /sales", h.GetSalesHandler)
	mux.HandleFunc("POST /sales", h.CreateSaleHandler)

	log.Println("Server starting on http://localhost:8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
