package loyalty

import (
	"context"
	"errors"
	"time"

	repo "github.com/Galish/loyalty-system/internal/repository"
	"github.com/ShiraazMoollatjie/goluhn"
)

type Order struct {
	ID         string `json:"number"`
	Status     Status `json:"status"`
	Accrual    uint   `json:"accrual"`
	UploadedAt string `json:"uploaded_at"`
	User       string `json:"user_id"`
}

type Status string

const (
	StatusNew        = Status("NEW")
	StatusRegistered = Status("REGISTERED")
	StatusProcessing = Status("PROCESSING")
	StatusInvalid    = Status("INVALID")
	StatusProcessed  = Status("PROCESSED")

	TimeLayout = time.RFC3339 //"2006-01-02T15:04:05-07:00"
)

var (
	ErrInvalidOrderID = errors.New("invalid order number value")
)

func (s *LoyaltyService) AddOrder(ctx context.Context, id, user string) (*Order, error) {
	if err := goluhn.Validate(id); err != nil {
		return nil, ErrInvalidOrderID
	}

	order := repo.Order{
		ID:         id,
		Status:     string(StatusNew),
		UploadedAt: time.Now(),
		User:       user,
	}

	if err := s.repo.CreateOrder(ctx, &order); err != nil {
		return nil, err
	}

	return &Order{
		ID:         order.ID,
		Status:     Status(order.Status),
		UploadedAt: order.UploadedAt.Format(TimeLayout),
		User:       user,
	}, nil
}

func (s *LoyaltyService) GetOrders(ctx context.Context, userID string) ([]*Order, error) {
	orders, err := s.repo.GetUserOrders(ctx, userID)
	if err != nil {
		return []*Order{}, nil
	}

	userOrders := []*Order{}

	for _, order := range orders {
		userOrders = append(
			userOrders,
			&Order{
				ID:         order.ID,
				Accrual:    order.Accrual,
				Status:     Status(order.Status),
				UploadedAt: order.UploadedAt.Format(TimeLayout),
				User:       order.User,
			},
		)
	}

	return userOrders, nil
}
