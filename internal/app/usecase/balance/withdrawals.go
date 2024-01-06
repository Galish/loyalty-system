package balance

import (
	"context"

	"github.com/Galish/loyalty-system/internal/app/entity"
)

func (uc *balanceUseCase) Withdrawals(ctx context.Context, user string) ([]*entity.Withdrawal, error) {
	withdrawals, err := uc.repo.Withdrawals(ctx, user)
	if err != nil {
		return []*entity.Withdrawal{}, err
	}

	return withdrawals, nil
}
