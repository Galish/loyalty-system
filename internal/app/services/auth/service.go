package auth

import (
	"context"

	"github.com/Galish/loyalty-system/internal/app/entity"
	"github.com/Galish/loyalty-system/internal/app/repository"
)

type AuthManager interface {
	Register(context.Context, string, string) (string, error)
	Authenticate(context.Context, string, string) (string, error)
	GenerateToken(*entity.User) (string, error)
	ParseToken(tokenString string) (*JWTClaims, error)
}

type AuthService struct {
	repo      repository.UserRepository
	secretKey string
}

func NewService(repo repository.UserRepository, secretKey string) *AuthService {
	return &AuthService{
		repo:      repo,
		secretKey: secretKey,
	}
}
