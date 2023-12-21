package balance

import (
	"context"
	"time"

	"github.com/Galish/loyalty-system/internal/model"
)

func (s *BalanceService) GetBalance(ctx context.Context, user string) (*model.Balance, error) {
	balance, err := s.repo.UserBalance(ctx, user)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

func (s *BalanceService) Withdraw(ctx context.Context, withdrawal *model.Withdrawal) error {
	err := s.repo.Withdraw(
		ctx,
		&model.Withdrawal{
			Order:       withdrawal.Order,
			User:        withdrawal.User,
			Sum:         withdrawal.Sum,
			ProcessedAt: model.Time(time.Now()),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *BalanceService) Withdrawals(ctx context.Context, user string) ([]*model.Withdrawal, error) {
	withdrawals, err := s.repo.Withdrawals(ctx, user)
	if err != nil {
		return []*model.Withdrawal{}, err
	}

	return withdrawals, nil
}
