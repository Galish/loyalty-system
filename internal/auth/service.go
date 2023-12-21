package auth

import (
	"context"

	"github.com/Galish/loyalty-system/internal/model"
	"github.com/Galish/loyalty-system/internal/repository"
)

type AuthManager interface {
	Register(context.Context, string, string) (string, error)
	Authenticate(context.Context, string, string) (string, error)
	GenerateToken(*model.User) (string, error)
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
