package auth

import (
	"context"
	"errors"

	"github.com/Galish/loyalty-system/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

const (
	AuthCookieName = "auth"
	AuthHeaderName = "X-User"
)

type AuthService struct {
	repo repository.UserRepository
}

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func NewService(repo repository.UserRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (as *AuthService) Register(ctx context.Context, creds Credentials) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user, err := as.repo.CreateUser(ctx, creds.Login, string(bytes))
	if err != nil {
		return "", err
	}

	token, err := GenerateToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (as *AuthService) Authenticate(ctx context.Context, creds Credentials) (string, error) {
	user, err := as.repo.GetUserByLogin(ctx, creds.Login)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		return "", errors.New("incorrect login/password pair")
	}

	token, err := GenerateToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}
