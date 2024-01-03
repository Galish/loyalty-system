package order

import (
	"context"
	"time"

	"github.com/Galish/loyalty-system/internal/app/entity"
)

func (s *OrderService) AddOrder(ctx context.Context, order entity.Order) error {
	if err := order.Validate(); err != nil {
		return err
	}

	order.Status = entity.StatusNew
	order.UploadedAt = entity.Time(time.Now())

	if err := s.repo.CreateOrder(ctx, &order); err != nil {
		return err
	}

	return nil
}
