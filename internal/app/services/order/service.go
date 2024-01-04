package order

import (
	"context"

	"github.com/Galish/loyalty-system/internal/app/entity"
	"github.com/Galish/loyalty-system/internal/app/repository"
)

type OrderManager interface {
	AddOrder(context.Context, entity.Order) error
	GetOrders(context.Context, string) ([]*entity.Order, error)
}

type orderService struct {
	repo repository.OrderRepository
}

func New(repo repository.OrderRepository) *orderService {
	return &orderService{
		repo: repo,
	}
}
