package user

import (
	"context"

	"github.com/Galish/loyalty-system/internal/auth"
	"golang.org/x/crypto/bcrypt"
)

func (uc *userUseCase) Authenticate(ctx context.Context, username, password string) (string, error) {
	user, err := uc.repo.GetUserByLogin(ctx, username)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", ErrIncorrectLoginPassword
	}

	token, err := auth.GenerateToken(uc.secretKey, user)
	if err != nil {
		return "", err
	}

	return token, nil
}
