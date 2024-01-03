package order

import (
	"context"
	"time"

	"github.com/Galish/loyalty-system/internal/entity"
)

func (s *OrderService) AddOrder(ctx context.Context, order entity.Order) error {
	order.Status = entity.StatusNew
	order.UploadedAt = entity.Time(time.Now())

	if err := s.repo.CreateOrder(ctx, &order); err != nil {
		return err
	}

	return nil
}

func (s *OrderService) GetOrders(ctx context.Context, userID string) ([]*entity.Order, error) {
	orders, err := s.repo.UserOrders(ctx, userID)
	if err != nil {
		return []*entity.Order{}, nil
	}

	return orders, nil
}
