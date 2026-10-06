package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"letsgo/internal/domain"
	"letsgo/internal/service"
)

type Handler struct {
	Service service.SalesService
}

func NewHandler(svc service.SalesService) *Handler {
	return &Handler{Service: svc}
}

func (h *Handler) GetSalesHandler(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")

	limit := 10
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit <= 0 {
			http.Error(w, "Invalid 'limit' parameter. Must be a positive integer.", http.StatusBadRequest)
			return
		}
		if parsedLimit > 100 {
			http.Error(w, "Limit cannot exceed 100.", http.StatusBadRequest)
			return
		}
		limit = parsedLimit
	}

	sales, err := h.Service.GetAll(limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var response []GetSaleResponseDTO
	for _, row := range sales {
		response = append(response, ToDTO(row))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) CreateSaleHandler(w http.ResponseWriter, r *http.Request) {
	var request PostSaleRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON request body.", http.StatusBadRequest)
		return
	}

	saleDate, err := time.Parse(time.RFC3339, request.SaleDate)
	if err != nil {
		http.Error(w, "Invalid 'sale_date'. Must be an RFC3339 timestamp.", http.StatusBadRequest)
		return
	}

	sale := domain.Sale{
		SaleDate:      saleDate,
		Category:      request.Category,
		ProductName:   request.ProductName,
		Quantity:      request.Quantity,
		TotalPrice:    request.TotalPrice,
		CustomerID:    request.CustomerID,
		Region:        request.Region,
		PaymentMethod: request.PaymentMethod,
	}
	if err := h.Service.Save(sale); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
