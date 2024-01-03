package psql

import (
	"context"

	"github.com/Galish/loyalty-system/internal/entity"
)

func (s *psqlStore) UpdateOrder(ctx context.Context, order *entity.Order) error {
	_, err := s.db.ExecContext(
		ctx,
		`
			UPDATE orders
			SET accrual = $1, status = $2
			WHERE uuid = $3
		`,
		order.Accrual,
		order.Status,
		order.ID,
	)
	if err != nil {
		return err
	}

	return nil
}
