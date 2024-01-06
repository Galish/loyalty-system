package balance

import (
	"context"

	"github.com/Galish/loyalty-system/internal/app/entity"
)

func (uc *balanceUseCase) GetBalance(ctx context.Context, user string) (*entity.Balance, error) {
	balance, err := uc.repo.UserBalance(ctx, user)
	if err != nil {
		return nil, err
	}

	return balance, nil
}
