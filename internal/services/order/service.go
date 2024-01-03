package order

import (
	"context"

	"github.com/Galish/loyalty-system/internal/entity"
	"github.com/Galish/loyalty-system/internal/repository"
)

type OrderManager interface {
	AddOrder(context.Context, entity.Order) error
	GetOrders(context.Context, string) ([]*entity.Order, error)
}

type OrderService struct {
	repo repository.OrderRepository
}

func NewService(repo repository.OrderRepository) *OrderService {
	return &OrderService{
		repo: repo,
	}
}
