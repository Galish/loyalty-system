package auth

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

const (
	AuthCookieName = "auth"
	AuthHeaderName = "X-User"
)

func (as *AuthService) Register(ctx context.Context, username, password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user, err := as.repo.CreateUser(ctx, username, string(bytes))
	if err != nil {
		return "", err
	}

	token, err := as.GenerateToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (as *AuthService) Authenticate(ctx context.Context, username, password string) (string, error) {
	user, err := as.repo.GetUserByLogin(ctx, username)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", ErrIncorrectLoginPassword
	}

	token, err := as.GenerateToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}
