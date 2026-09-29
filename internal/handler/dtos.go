package handler

import (
	"letsgo/internal/domain"
	"time"
)

type GetSaleResponseDTO struct {
	ID       int     `json:"id"`
	Product  string  `json:"product_name"`
	Revenue  float64 `json:"revenue"`
	SaleDate string  `json:"sale_date"`
}

type PostSaleRequestDTO struct {
	SaleDate      string  `json:"sale_date"`
	Category      string  `json:"category"`
	ProductName   string  `json:"product_name"`
	Quantity      int     `json:"quantity"`
	TotalPrice    float64 `json:"total_price"`
	CustomerID    int     `json:"customer_id"`
	Region        string  `json:"region"`
	PaymentMethod string  `json:"payment_method"`
}

func ToDTO(sale domain.Sale) GetSaleResponseDTO {
	return GetSaleResponseDTO{
		ID:       sale.ID,
		Product:  sale.ProductName,
		Revenue:  sale.TotalPrice,
		SaleDate: sale.SaleDate.Format(time.RFC3339),
	}
}
