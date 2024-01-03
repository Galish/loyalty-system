package balance

import (
	"context"
	"time"

	"github.com/Galish/loyalty-system/internal/entity"
)

func (s *BalanceService) GetBalance(ctx context.Context, user string) (*entity.Balance, error) {
	balance, err := s.repo.UserBalance(ctx, user)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

func (s *BalanceService) Withdraw(ctx context.Context, withdrawal *entity.Withdrawal) error {
	err := s.repo.Withdraw(
		ctx,
		&entity.Withdrawal{
			Order:       withdrawal.Order,
			User:        withdrawal.User,
			Sum:         withdrawal.Sum,
			ProcessedAt: entity.Time(time.Now()),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *BalanceService) Withdrawals(ctx context.Context, user string) ([]*entity.Withdrawal, error) {
	withdrawals, err := s.repo.Withdrawals(ctx, user)
	if err != nil {
		return []*entity.Withdrawal{}, err
	}

	return withdrawals, nil
}
