package loyalty

import (
	"context"
	"errors"
	"time"

	"github.com/Galish/loyalty-system/internal/repository"
)

type Balance struct {
	Current   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}

type Withdraw struct {
	Order       string    `json:"order"`
	Sum         float32   `json:"sum"`
	User        string    `json:"user,omitempty"`
	ProcessedAt time.Time `json:"processed_at"`
}

var ErrInsufficientFunds = errors.New("insufficient funds in the account")

func (s *LoyaltyService) GetBalance(ctx context.Context, user string) (*Balance, error) {
	balance, err := s.repo.GetUserBalance(ctx, user)
	if err != nil {
		return nil, err
	}

	return &Balance{
		Current:   balance.Current,
		Withdrawn: balance.Withdrawn,
	}, nil
}

func (s *LoyaltyService) Withdraw(ctx context.Context, withdrawn *Withdraw) error {
	if !s.ValidateOrderNumber(withdrawn.Order) {
		return ErrInvalidOrderNumber
	}

	balance, err := s.repo.GetUserBalance(ctx, withdrawn.User)
	if err != nil {
		return err
	}
	if balance.Current < withdrawn.Sum {
		return ErrInsufficientFunds
	}

	if err := s.repo.UpdateBalance(ctx, withdrawn.User, -withdrawn.Sum); err != nil {
		return err
	}

	return s.repo.CreateWithdraw(ctx, &repository.Withdraw{
		Order:       withdrawn.Order,
		User:        withdrawn.User,
		Sum:         withdrawn.Sum,
		ProcessedAt: time.Now(),
	})
}

func (s *LoyaltyService) GetWithdrawals(ctx context.Context, user string) ([]*Withdraw, error) {
	withdrawals, err := s.repo.GetWithdrawals(ctx, user)
	if err != nil {
		return []*Withdraw{}, err
	}

	var results []*Withdraw
	for _, w := range withdrawals {
		results = append(
			results,
			&Withdraw{
				Order:       w.Order,
				Sum:         w.Sum,
				ProcessedAt: w.ProcessedAt,
			},
		)
	}

	return results, nil
}
