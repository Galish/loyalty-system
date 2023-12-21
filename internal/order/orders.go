package order

import (
	"context"
	"time"

	"github.com/Galish/loyalty-system/internal/model"
)

func (s *OrderService) AddOrder(ctx context.Context, order model.Order) error {
	if !order.ID.IsValid() {
		return ErrIncorrectOrderNumber
	}

	order.Status = model.StatusNew
	order.UploadedAt = time.Now()

	if err := s.repo.CreateOrder(ctx, &order); err != nil {
		return err
	}

	return nil
}

func (s *OrderService) GetOrders(ctx context.Context, userID string) ([]*model.Order, error) {
	orders, err := s.repo.UserOrders(ctx, userID)
	if err != nil {
		return []*model.Order{}, nil
	}

	return orders, nil
}
