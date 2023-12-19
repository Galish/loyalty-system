package loyalty

import (
	"context"
	"errors"
	"time"

	repo "github.com/Galish/loyalty-system/internal/repository"
)

type Order struct {
	ID         OrderNumber `json:"number"`
	Status     Status      `json:"status"`
	Accrual    float32     `json:"accrual"`
	UploadedAt string      `json:"uploaded_at"`
	User       string      `json:"user_id,omitempty"`
}

var ErrIncorrectOrderNumber = errors.New("invalid order number value")

func (s *LoyaltyService) AddOrder(ctx context.Context, order *Order) error {
	if !order.ID.isValid() {
		return ErrIncorrectOrderNumber
	}

	repoOrder := repo.Order{
		ID:         order.ID.String(),
		Status:     string(StatusNew),
		UploadedAt: time.Now(),
		User:       order.User,
	}

	if err := s.repo.CreateOrder(ctx, &repoOrder); err != nil {
		return err
	}

	s.orderCh <- order

	return nil
}

func (s *LoyaltyService) GetOrders(ctx context.Context, userID string) ([]*Order, error) {
	orders, err := s.repo.UserOrders(ctx, userID)
	if err != nil {
		return []*Order{}, nil
	}

	userOrders := []*Order{}

	for _, order := range orders {
		userOrders = append(
			userOrders,
			&Order{
				ID:         OrderNumber(order.ID),
				Accrual:    order.Accrual,
				Status:     Status(order.Status),
				UploadedAt: order.UploadedAt.Format(TimeLayout),
			},
		)
	}

	return userOrders, nil
}
