package repository

import (
	"context"
	"errors"
)

var (
	ErrConflict = errors.New("user already exists")
	ErrNotFound = errors.New("user not found")
)

type UserRepository interface {
	Create(context.Context, string, string) (*User, error)
	GetByLogin(context.Context, string) (*User, error)
}

type LoyaltyRepository interface {
}

type User struct {
	ID       string `json:"uuid"`
	Login    string `json:"login"`
	Password string `json:"password"`
	IsActive bool   `json:"is_active"`
}
