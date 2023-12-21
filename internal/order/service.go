package order

import (
	"github.com/Galish/loyalty-system/internal/repository"
	repo "github.com/Galish/loyalty-system/internal/repository"
)

type OrderService struct {
	repo repository.OrderRepository
}

func NewService(repo repo.OrderRepository) *OrderService {
	return &OrderService{
		repo: repo,
	}
}
