package auth

import (
	"context"
	"loyalty-system/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo repository.UserRepository
}

func NewService(repo repository.UserRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (as *AuthService) Register(ctx context.Context, login, password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}

	return as.repo.Create(ctx, login, string(bytes))
}

func (as *AuthService) Authenticate(ctx context.Context, login, password string) error {
	user, err := as.repo.GetByLogin(ctx, login)
	if err != nil {
		return nil
	}

	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}
