package loyalty

import (
	"context"
	"errors"
	"time"

	"github.com/Galish/loyalty-system/internal/model"
)

const TimeLayout = "2006-01-02T15:04:05-07:00"

var ErrIncorrectOrderNumber = errors.New("invalid order number value")

func (s *LoyaltyService) AddOrder(ctx context.Context, order model.Order) error {
	if !order.ID.IsValid() {
		return ErrIncorrectOrderNumber
	}

	order.Status = model.StatusNew
	order.UploadedAt = time.Now()

	if err := s.repo.CreateOrder(ctx, &order); err != nil {
		return err
	}

	go func() {
		s.orderCh <- &order
	}()

	return nil
}

func (s *LoyaltyService) GetOrders(ctx context.Context, userID string) ([]*model.Order, error) {
	orders, err := s.repo.UserOrders(ctx, userID)
	if err != nil {
		return []*model.Order{}, nil
	}

	return orders, nil
}