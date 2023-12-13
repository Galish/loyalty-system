package psql

import (
	"context"

	repo "github.com/Galish/loyalty-system/internal/repository"
)

func (s *psqlStore) GetWithdrawals(ctx context.Context, user string) ([]*repo.Withdraw, error) {
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
		return []*repo.Withdraw{}, err
	}

	defer rows.Close()

	var withdrawals []*repo.Withdraw

	for rows.Next() {
		var withdraw repo.Withdraw

		if err := rows.Scan(
			&withdraw.Order,
			&withdraw.User,
			&withdraw.Sum,
			&withdraw.ProcessedAt,
		); err != nil {
			return []*repo.Withdraw{}, err
		}

		withdrawals = append(withdrawals, &withdraw)
	}

	if err := rows.Err(); err != nil {
		return []*repo.Withdraw{}, err
	}

	return withdrawals, nil
}
