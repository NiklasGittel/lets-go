package repository

import (
	"database/sql"
	"errors"
	"letsgo/internal/domain"
)

type SQLiteRepo struct {
	db *sql.DB
}

// Constructor function
func NewSQLiteRepo(db *sql.DB) SaleRepository {
	return &SQLiteRepo{db: db}
}

func (r *SQLiteRepo) GetByID(id int) (*domain.Sale, error) {
	query := `
		SELECT transaction_id, sale_date, category, product_name, quantity, total_price, customer_id, region, payment_method 
		FROM sales WHERE transaction_id = ?
	`
	row := r.db.QueryRow(query, id)

	var s domain.Sale
	err := row.Scan(
		&s.ID, &s.SaleDate, &s.Category, &s.ProductName,
		&s.Quantity, &s.TotalPrice, &s.CustomerID, &s.Region, &s.PaymentMethod,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *SQLiteRepo) GetAll(limit int) ([]domain.Sale, error) {
	query := `
        SELECT transaction_id, sale_date, category, product_name, quantity, total_price, customer_id, region, payment_method 
        FROM sales
    `

	var rows *sql.Rows
	var err error

	if limit > 0 {
		query += " LIMIT ?"
		rows, err = r.db.Query(query, limit)
	} else {
		rows, err = r.db.Query(query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sales []domain.Sale
	for rows.Next() {
		var s domain.Sale
		if err := rows.Scan(
			&s.ID, &s.SaleDate, &s.Category, &s.ProductName,
			&s.Quantity, &s.TotalPrice, &s.CustomerID, &s.Region, &s.PaymentMethod,
		); err != nil {
			return nil, err
		}
		sales = append(sales, s)
	}
	return sales, nil
}

func (r *SQLiteRepo) Save(s domain.Sale) error {
	query := `
		INSERT INTO sales (sale_date, category, product_name, quantity, total_price, customer_id, region, payment_method)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(
		query, s.SaleDate, s.Category, s.ProductName,
		s.Quantity, s.TotalPrice, s.CustomerID, s.Region, s.PaymentMethod,
	)
	return err
}
