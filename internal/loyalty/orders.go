package loyalty

import (
	"context"
	"errors"
	"time"

	repo "github.com/Galish/loyalty-system/internal/repository"
	"github.com/ShiraazMoollatjie/goluhn"
)

type Order struct {
	ID         OrderNumber `json:"number"`
	Status     Status      `json:"status"`
	Accrual    float32     `json:"accrual"`
	UploadedAt string      `json:"uploaded_at"`
	User       string      `json:"user_id,omitempty"`
}

type OrderNumber string

type Status string

const (
	StatusNew        = Status("NEW")
	StatusRegistered = Status("REGISTERED")
	StatusProcessing = Status("PROCESSING")
	StatusInvalid    = Status("INVALID")
	StatusProcessed  = Status("PROCESSED")

	TimeLayout = "2006-01-02T15:04:05-07:00"
)

var ErrIncorrectOrderNumber = errors.New("invalid order number value")

func (s *LoyaltyService) AddOrder(ctx context.Context, order *Order) error {
	if !order.ID.isValid() {
		return ErrIncorrectOrderNumber
	}

	repoOrder := repo.Order{
		ID:         order.ID.String(),
		Status:     string(StatusNew),
		UploadedAt: time.Now().Round(time.Microsecond),
		User:       order.User,
	}

	if err := s.repo.CreateOrder(ctx, &repoOrder); err != nil {
		return err
	}

	go func() {
		s.orderCh <- order
	}()

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

func (num OrderNumber) isValid() bool {
	if num.String() == "" {
		return false
	}

	if err := goluhn.Validate(num.String()); err != nil {
		return false
	}

	return true
}

func (num OrderNumber) String() string {
	return string(num)
}
