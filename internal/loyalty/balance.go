package loyalty

import (
	"context"
	"time"

	"github.com/Galish/loyalty-system/internal/model"
	repo "github.com/Galish/loyalty-system/internal/repository"
)

type Balance struct {
	Current   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}

type Withdrawal struct {
	Order       model.OrderNumber `json:"order"`
	Sum         float32           `json:"sum"`
	User        string            `json:"user,omitempty"`
	ProcessedAt string            `json:"processed_at"`
}

func (s *LoyaltyService) GetBalance(ctx context.Context, user string) (*Balance, error) {
	balance, err := s.repo.UserBalance(ctx, user)
	if err != nil {
		return nil, err
	}

	return &Balance{
		Current:   balance.Current,
		Withdrawn: balance.Withdrawn,
	}, nil
}

func (s *LoyaltyService) Withdraw(ctx context.Context, withdrawn *Withdrawal) error {
	if !withdrawn.Order.IsValid() {
		return ErrIncorrectOrderNumber
	}

	err := s.repo.Withdraw(
		ctx,
		&repo.Withdrawal{
			Order:       withdrawn.Order.String(),
			User:        withdrawn.User,
			Sum:         withdrawn.Sum,
			ProcessedAt: time.Now(),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *LoyaltyService) Withdrawals(ctx context.Context, user string) ([]*Withdrawal, error) {
	withdrawals, err := s.repo.Withdrawals(ctx, user)
	if err != nil {
		return []*Withdrawal{}, err
	}

	var results []*Withdrawal
	for _, w := range withdrawals {
		results = append(
			results,
			&Withdrawal{
				Order:       model.OrderNumber(w.Order),
				Sum:         w.Sum,
				ProcessedAt: w.ProcessedAt.Format(TimeLayout),
			},
		)
	}

	return results, nil
}
