package balance

import (
	"context"
	"time"

	"github.com/Galish/loyalty-system/internal/app/entity"
)

func (s *balanceService) Withdraw(ctx context.Context, withdrawal *entity.Withdrawal) error {
	if err := withdrawal.Validate(); err != nil {
		return err
	}

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
