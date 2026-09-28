package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"letsgo/internal/repository"
)

type Handler struct {
	Repo repository.SaleRepository
}

func NewHandler(repo repository.SaleRepository) *Handler {
	return &Handler{Repo: repo}
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

	sales, err := h.Repo.GetAll(limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var response []SaleResponseDTO
	for _, row := range sales {
		response = append(response, ToDTO(row))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
