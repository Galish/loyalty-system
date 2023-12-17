package psql

import (
	"context"

	repo "github.com/Galish/loyalty-system/internal/repository"
)

func (s *psqlStore) Withdrawals(ctx context.Context, user string) ([]*repo.Withdrawal, error) {
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
		return []*repo.Withdrawal{}, err
	}

	defer rows.Close()

	var withdrawals []*repo.Withdrawal

	for rows.Next() {
		var withdraw repo.Withdrawal

		if err := rows.Scan(
			&withdraw.Order,
			&withdraw.User,
			&withdraw.Sum,
			&withdraw.ProcessedAt,
		); err != nil {
			return []*repo.Withdrawal{}, err
		}

		withdrawals = append(withdrawals, &withdraw)
	}

	if err := rows.Err(); err != nil {
		return []*repo.Withdrawal{}, err
	}

	return withdrawals, nil
}
