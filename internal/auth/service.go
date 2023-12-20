package auth

import (
	"github.com/Galish/loyalty-system/internal/repository"
)

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
