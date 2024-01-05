package balance

import (
	"context"
	"errors"
	"time"

	"github.com/Galish/loyalty-system/internal/app/entity"
)

var ErrInvalidOrderNumber = errors.New("invalid withdrawal order number")

func (uc *balanceUseCase) Withdraw(ctx context.Context, withdrawal *entity.Withdrawal) error {
	if !withdrawal.IsValid() {
		return ErrInvalidOrderNumber
	}

	err := uc.repo.Withdraw(
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
