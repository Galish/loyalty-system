package psql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Galish/loyalty-system/internal/model"
	repo "github.com/Galish/loyalty-system/internal/repository"
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
	err := row.Scan(
		&balance.User,
		&balance.Current,
		&balance.Withdrawn,
		&balance.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repo.ErrNothingFound
	}
	if err != nil {
		return nil, err
	}

	return &balance, nil
}
