package psql

import (
	"context"

	"github.com/Galish/loyalty-system/internal/model"
)

func (s *psqlStore) UserBalance(ctx context.Context, user string) (*model.Balance, error) {
	row := s.db.QueryRowContext(
		ctx,
		`
			SELECT user_id, current, withdrawn, updated_at
			FROM balance
			WHERE user_id = $1
		`,
		user,
	)

	var balance model.Balance
	if err := row.Scan(
		&balance.User,
		&balance.Current,
		&balance.Withdrawn,
		&balance.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &balance, nil
}
