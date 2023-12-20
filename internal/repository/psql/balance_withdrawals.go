package psql

import (
	"context"

	"github.com/Galish/loyalty-system/internal/model"
)

func (s *psqlStore) Withdrawals(ctx context.Context, user string) ([]*model.Withdrawal, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`
			SELECT order_id, user_id, sum, processed_at
			FROM withdrawals
			WHERE user_id = $1
		`,
		user,
	)
	if err != nil {
		return []*model.Withdrawal{}, err
	}

	defer rows.Close()

	var withdrawals []*model.Withdrawal

	for rows.Next() {
		var withdraw model.Withdrawal

		if err := rows.Scan(
			&withdraw.Order,
			&withdraw.User,
			&withdraw.Sum,
			&withdraw.ProcessedAt,
		); err != nil {
			return []*model.Withdrawal{}, err
		}

		withdrawals = append(withdrawals, &withdraw)
	}

	if err := rows.Err(); err != nil {
		return []*model.Withdrawal{}, err
	}

	return withdrawals, nil
}
