package psql

import (
	"context"
	"database/sql"
	"errors"

	repo "github.com/Galish/loyalty-system/internal/repository"
)

func (s *psqlStore) GetUserByLogin(ctx context.Context, login string) (*repo.User, error) {
	row := s.db.QueryRowContext(
		ctx,
		`SELECT uuid, login, password, is_active FROM users WHERE login=$1;`,
		login,
	)

	var user repo.User
	err := row.Scan(
		&user.ID,
		&user.Login,
		&user.Password,
		&user.IsActive,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, repo.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}
