package user

import (
	"github.com/Galish/loyalty-system/internal/app/repository"
)

type authService struct {
	repo      repository.UserRepository
	secretKey string
}

func New(repo repository.UserRepository, secretKey string) *authService {
	return &authService{
		repo:      repo,
		secretKey: secretKey,
	}
}
