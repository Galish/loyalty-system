package balance

import (
	"context"

	"github.com/Galish/loyalty-system/internal/app/entity"
)

func (s *balanceService) Withdrawals(ctx context.Context, user string) ([]*entity.Withdrawal, error) {
	withdrawals, err := s.repo.Withdrawals(ctx, user)
	if err != nil {
		return []*entity.Withdrawal{}, err
	}

	return withdrawals, nil
}
