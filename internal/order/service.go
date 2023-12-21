package order

import (
	"context"

	"github.com/Galish/loyalty-system/internal/model"
	"github.com/Galish/loyalty-system/internal/repository"
)

type OrderManager interface {
	AddOrder(context.Context, model.Order) error
	GetOrders(context.Context, string) ([]*model.Order, error)
}

type OrderService struct {
	repo repository.OrderRepository
}

func NewService(repo repository.OrderRepository) *OrderService {
	return &OrderService{
		repo: repo,
	}
}
