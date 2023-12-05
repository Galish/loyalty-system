package userRepository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func (s *psqlStore) Create(ctx context.Context, login, password string) error {
	_, err := s.db.ExecContext(
		ctx,
		`
			INSERT INTO users (login, password)
			VALUES ($1, $2);
		`,
		login,
		password,
	)

	var pgErr *pgconn.PgError
	if err != nil && errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrConflict
	}

	return nil
}
