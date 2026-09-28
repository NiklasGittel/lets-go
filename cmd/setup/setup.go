package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "modernc.org/sqlite"
)

const (
	dbName     = "local_sales.db"
	numRecords = 5000
)

type Product struct {
	Name  string
	Price float64
}

type SaleRecord struct {
	SaleDate      string
	Category      string
	ProductName   string
	Quantity      int
	UnitPrice     float64
	TotalPrice    float64
	CustomerID    int
	Region        string
	PaymentMethod string
}

var categoriesAndProducts = map[string][]Product{
	"Electronics": {{"Laptop", 999.99}, {"Smartphone", 699.50}, {"Headphones", 149.99}, {"Monitor", 249.00}},
	"Apparel":     {{"T-Shirt", 19.99}, {"Jeans", 49.99}, {"Jacket", 89.50}, {"Sneakers", 75.00}},
}
var regions = []string{"North America", "Europe", "Asia Pacific"}
var paymentMethods = []string{"Credit Card", "PayPal", "Debit Card"}

// 1. Function to generate mock data
func generateMockData(count int) []SaleRecord {
	var records []SaleRecord
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	startDate := time.Now().AddDate(-1, 0, 0)

	categories := make([]string, 0, len(categoriesAndProducts))
	for cat := range categoriesAndProducts {
		categories = append(categories, cat)
	}

	for i := 0; i < count; i++ {
		randomDays := r.Intn(366)
		saleDate := startDate.AddDate(0, 0, randomDays).Format("2006-01-02 15:04:05")

		category := categories[r.Intn(len(categories))]
		products := categoriesAndProducts[category]
		product := products[r.Intn(len(products))]

		quantity := r.Intn(5) + 1
		totalPrice := float64(quantity) * product.Price

		records = append(records, SaleRecord{
			SaleDate:      saleDate,
			Category:      category,
			ProductName:   product.Name,
			Quantity:      quantity,
			UnitPrice:     product.Price,
			TotalPrice:    totalPrice,
			CustomerID:    r.Intn(9000) + 1000,
			Region:        regions[r.Intn(len(regions))],
			PaymentMethod: paymentMethods[r.Intn(len(paymentMethods))],
		})
	}
	return records
}

// 2. Function to setup database and insert data
func setupDatabase(records []SaleRecord) {
	db, err := sql.Open("sqlite", dbName)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS sales (
		transaction_id INTEGER PRIMARY KEY AUTOINCREMENT,
		sale_date DATETIME,
		category TEXT,
		product_name TEXT,
		quantity INTEGER,
		unit_price REAL,
		total_price REAL,
		customer_id INTEGER,
		region TEXT,
		payment_method TEXT
	);`

	if _, err := db.Exec(createTableSQL); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	if _, err := db.Exec("DELETE FROM sales"); err != nil {
		log.Fatalf("Failed to clear table: %v", err)
	}

	// Bulk insert using a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}

	stmt, err := tx.Prepare(`
		INSERT INTO sales (
			sale_date, category, product_name, quantity, 
			unit_price, total_price, customer_id, region, payment_method
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	for _, rec := range records {
		_, err = stmt.Exec(
			rec.SaleDate, rec.Category, rec.ProductName, rec.Quantity,
			rec.UnitPrice, rec.TotalPrice, rec.CustomerID, rec.Region, rec.PaymentMethod,
		)
		if err != nil {
			tx.Rollback()
			log.Fatalf("Failed to insert row: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	fmt.Println("Database setup and insertion complete!")
}

func main() {
	fmt.Printf("Generating %d records...\n", numRecords)
	data := generateMockData(numRecords)

	fmt.Println("Saving to database...")
	setupDatabase(data)
}
