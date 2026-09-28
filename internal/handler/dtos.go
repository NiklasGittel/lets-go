package handler

import (
	"letsgo/internal/domain"
	"time"
)

type SaleResponseDTO struct {
	ID       int     `json:"id"`
	Product  string  `json:"product_name"`
	Revenue  float64 `json:"revenue"`
	SaleDate string  `json:"sale_date"`
}

func ToDTO(sale domain.Sale) SaleResponseDTO {
	return SaleResponseDTO{
		ID:       sale.ID,
		Product:  sale.ProductName,
		Revenue:  sale.TotalPrice,
		SaleDate: sale.SaleDate.Format(time.RFC3339),
	}
}
