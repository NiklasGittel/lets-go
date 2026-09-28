package domain

import "time"

type Sale struct {
	ID            int
	SaleDate      time.Time
	Category      string
	ProductName   string
	Quantity      int
	TotalPrice    float64
	CustomerID    int
	Region        string
	PaymentMethod string
}
