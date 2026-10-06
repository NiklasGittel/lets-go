package service

import (
	"slices"
	"sync"
	"time"

	"letsgo/internal/domain"
	"letsgo/internal/repository"
)

type SalesService interface {
	GetByID(id int) (*domain.Sale, error)
	GetAll(limit int) ([]domain.Sale, error)
	Save(sale domain.Sale) error
}

type cacheEntry[T any] struct {
	value     T
	expiresAt time.Time
}

type salesService struct {
	repo repository.SaleRepository
	ttl  time.Duration

	mu      sync.RWMutex
	byID    map[int]cacheEntry[domain.Sale]
	listing map[int]cacheEntry[[]domain.Sale]
}

func NewSalesService(repo repository.SaleRepository, ttl time.Duration) SalesService {
	return &salesService{
		repo:    repo,
		ttl:     ttl,
		byID:    make(map[int]cacheEntry[domain.Sale]),
		listing: make(map[int]cacheEntry[[]domain.Sale]),
	}
}

func (s *salesService) GetByID(id int) (*domain.Sale, error) {
	s.mu.RLock()
	e, ok := s.byID[id]
	s.mu.RUnlock()
	if ok && time.Now().Before(e.expiresAt) {
		sale := e.value
		return &sale, nil
	}

	sale, err := s.repo.GetByID(id)
	if err != nil || sale == nil {
		return sale, err
	}

	s.mu.Lock()
	s.byID[id] = cacheEntry[domain.Sale]{value: *sale, expiresAt: time.Now().Add(s.ttl)}
	s.mu.Unlock()
	return sale, nil
}

func (s *salesService) GetAll(limit int) ([]domain.Sale, error) {
	s.mu.RLock()
	e, ok := s.listing[limit]
	s.mu.RUnlock()
	if ok && time.Now().Before(e.expiresAt) {
		return slices.Clone(e.value), nil
	}

	sales, err := s.repo.GetAll(limit)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	s.listing[limit] = cacheEntry[[]domain.Sale]{value: slices.Clone(sales), expiresAt: time.Now().Add(s.ttl)}
	s.mu.Unlock()
	return sales, nil
}

func (s *salesService) Save(sale domain.Sale) error {
	if err := s.repo.Save(sale); err != nil {
		return err
	}
	s.mu.Lock()
	clear(s.byID)
	clear(s.listing)
	s.mu.Unlock()
	return nil
}
