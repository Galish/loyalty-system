package psql

import (
	"context"

	repo "github.com/Galish/loyalty-system/internal/repository"
)

func (s *psqlStore) CreateWithdraw(ctx context.Context, withdraw *repo.Withdraw) error {
	_, err := s.db.ExecContext(
		ctx,
		`
			INSERT INTO withdrawals (order_id, user_id, sum, processed_at)
			VALUES ($1, $2, $3, $4)
		`,
		withdraw.Order,
		withdraw.User,
		withdraw.Sum,
		withdraw.ProcessedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
