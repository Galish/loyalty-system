package loyalty

import (
	"context"
	"errors"
	"time"

	repo "github.com/Galish/loyalty-system/internal/repository"
)

type Balance struct {
	Current   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}

type Withdraw struct {
	Order       string  `json:"order"`
	Sum         float32 `json:"sum"`
	User        string  `json:"user,omitempty"`
	ProcessedAt string  `json:"processed_at"`
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

	err := s.repo.Withdraw(
		ctx,
		&repo.Withdraw{
			Order:       withdrawn.Order,
			User:        withdrawn.User,
			Sum:         withdrawn.Sum,
			ProcessedAt: time.Now(),
		},
	)
	if err != nil && errors.Is(err, repo.ErrInsufficientFunds) {
		return err
	}
	if err != nil {
		return err
	}

	return nil
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
				ProcessedAt: w.ProcessedAt.Format(TimeLayout),
			},
		)
	}

	return results, nil
}
