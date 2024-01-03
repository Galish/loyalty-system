package order

import (
	"github.com/Galish/loyalty-system/internal/app/repository"
)

// type OrderManager interface {
// 	AddOrder(context.Context, entity.Order) error
// 	GetOrders(context.Context, string) ([]*entity.Order, error)
// }

type OrderService struct {
	repo repository.OrderRepository
}

func NewService(repo repository.OrderRepository) *OrderService {
	return &OrderService{
		repo: repo,
	}
}
