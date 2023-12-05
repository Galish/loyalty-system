package userRepository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

func (s *psqlStore) GetByLogin(ctx context.Context, login string) (*User, error) {
	row := s.db.QueryRowContext(
		ctx,
		`SELECT * FROM users WHERE login=$1;`,
		login,
	)

	var user User
	err := row.Scan(
		&user.ID,
		&user.Login,
		&user.Password,
		&user.IsActive,
	)

	if err != nil && errors.Is(pgx.ErrNoRows, err) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}
