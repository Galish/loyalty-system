package psql

import (
	"context"
	"errors"

	repo "github.com/Galish/loyalty-system/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *psqlStore) Create(ctx context.Context, login, password string) (*repo.User, error) {
	user := repo.User{
		ID:       uuid.NewString(),
		Login:    login,
		Password: password,
	}
	_, err := s.db.ExecContext(
		ctx,
		`
			INSERT INTO users (uuid, login, password)
			VALUES ($1, $2, $3);
		`,
		user.ID,
		user.Login,
		user.Password,
	)

	var pgErr *pgconn.PgError
	if err != nil && errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return nil, repo.ErrConflict
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}
