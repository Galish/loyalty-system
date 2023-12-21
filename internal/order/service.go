package order

import (
	"github.com/Galish/loyalty-system/internal/repository"
)

type OrderService struct {
	repo repository.OrderRepository
}

func NewService(repo repository.OrderRepository) *OrderService {
	return &OrderService{
		repo: repo,
	}
}
