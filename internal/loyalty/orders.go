package loyalty

import (
	"context"
	"errors"
	"time"

	repo "github.com/Galish/loyalty-system/internal/repository"
	"github.com/ShiraazMoollatjie/goluhn"
)

type Order struct {
	ID         string  `json:"order"`
	Status     Status  `json:"status"`
	Accrual    float32 `json:"accrual"`
	UploadedAt string  `json:"uploaded_at"`
	User       string  `json:"user_id"`
}

type Status string

const (
	StatusNew        = Status("NEW")
	StatusRegistered = Status("REGISTERED")
	StatusProcessing = Status("PROCESSING")
	StatusInvalid    = Status("INVALID")
	StatusProcessed  = Status("PROCESSED")

	TimeLayout = time.RFC3339
)

var (
	ErrInvalidOrderNumber = errors.New("invalid order number value")
)

func (s *LoyaltyService) AddOrder(ctx context.Context, id, user string) (*Order, error) {
	if !s.ValidateOrderNumber(id) {
		return nil, ErrInvalidOrderNumber
	}

	repoOrder := repo.Order{
		ID:         id,
		Status:     string(StatusNew),
		UploadedAt: time.Now(),
		User:       user,
	}

	if err := s.repo.CreateOrder(ctx, &repoOrder); err != nil {
		return nil, err
	}

	order := Order{
		ID:         repoOrder.ID,
		Status:     Status(repoOrder.Status),
		UploadedAt: repoOrder.UploadedAt.Format(TimeLayout),
		User:       user,
	}

	go func() {
		s.orderCh <- &order
	}()

	return &order, nil
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

func (s *LoyaltyService) ValidateOrderNumber(id string) bool {
	if id == "" {
		return false
	}

	if err := goluhn.Validate(id); err != nil {
		return false
	}

	return true
}
