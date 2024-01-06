package user

import (
	"context"

	"github.com/Galish/loyalty-system/internal/auth"
	"golang.org/x/crypto/bcrypt"
)

func (uc *userUseCase) Register(ctx context.Context, username, password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user, err := uc.repo.CreateUser(ctx, username, string(bytes))
	if err != nil {
		return "", err
	}

	token, err := auth.GenerateToken(uc.secretKey, user)
	if err != nil {
		return "", err
	}

	return token, nil
}
