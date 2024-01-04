package order

import (
	"context"

	"github.com/Galish/loyalty-system/internal/app/entity"
)

func (s *orderService) GetOrders(ctx context.Context, userID string) ([]*entity.Order, error) {
	orders, err := s.repo.UserOrders(ctx, userID)
	if err != nil {
		return []*entity.Order{}, nil
	}

	return orders, nil
}
