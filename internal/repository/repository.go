package repository

import (
	"context"
	userRepo "loyalty-system/internal/repository/user"
)

type UserRepository interface {
	Create(context.Context, string, string) (*userRepo.User, error)
	GetByLogin(context.Context, string) (*userRepo.User, error)
}

type LoyaltyRepository interface {
}
