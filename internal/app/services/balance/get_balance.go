package balance

import (
	"context"

	"github.com/Galish/loyalty-system/internal/app/entity"
)

func (s *balanceService) GetBalance(ctx context.Context, user string) (*entity.Balance, error) {
	balance, err := s.repo.UserBalance(ctx, user)
	if err != nil {
		return nil, err
	}

	return balance, nil
}
