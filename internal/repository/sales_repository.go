package repository

import "letsgo/internal/domain"

type SaleRepository interface {
	GetByID(id int) (*domain.Sale, error)
	GetAll(limit int) ([]domain.Sale, error)
	Save(sale domain.Sale) error
}
