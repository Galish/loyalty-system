package psql

import (
	"context"
	"errors"

	repo "loyalty-system/internal/repository"

	"github.com/jackc/pgx/v5"
)

func (s *psqlStore) GetByLogin(ctx context.Context, login string) (*repo.User, error) {
	row := s.db.QueryRowContext(
		ctx,
		`SELECT * FROM users WHERE login=$1;`,
		login,
	)

	var user repo.User
	err := row.Scan(
		&user.ID,
		&user.Login,
		&user.Password,
		&user.IsActive,
	)

	if err != nil && errors.Is(pgx.ErrNoRows, err) {
		return nil, repo.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}
