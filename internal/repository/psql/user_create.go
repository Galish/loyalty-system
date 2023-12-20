package psql

import (
	"context"
	"errors"

	"github.com/Galish/loyalty-system/internal/model"
	repo "github.com/Galish/loyalty-system/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *psqlStore) CreateUser(ctx context.Context, username, password string) (*model.User, error) {
	user := model.User{
		ID:       uuid.NewString(),
		Login:    username,
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
		return nil, repo.ErrUserConflict
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}
