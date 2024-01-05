package user

import (
	"github.com/Galish/loyalty-system/internal/app/repository"
)

type userUseCase struct {
	repo      repository.UserRepository
	secretKey string
}

func New(repo repository.UserRepository, secretKey string) *userUseCase {
	return &userUseCase{
		repo:      repo,
		secretKey: secretKey,
	}
}
