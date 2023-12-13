package repository

import (
	"context"
	"errors"
	"time"
)

var (
	ErrOrderExists   = errors.New("order has already been added")
	ErrOrderConflict = errors.New("order has already been added by another user")
	ErrUserConflict  = errors.New("user already exists")
	ErrUserNotFound  = errors.New("user not found")
)

type UserRepository interface {
	CreateUser(context.Context, string, string) (*User, error)
	GetUserByLogin(context.Context, string) (*User, error)
}

type LoyaltyRepository interface {
	CreateOrder(context.Context, *Order) error
	GetUserOrders(context.Context, string) ([]*Order, error)
	UpdateOrder(context.Context, *Order) error
	UpdateBalance(context.Context, string, int) error
}

type User struct {
	ID       string `json:"uuid"`
	Login    string `json:"login"`
	Password string `json:"password"`
	IsActive bool   `json:"is_active"`
}

type Order struct {
	ID         string    `json:"uuid"`
	Status     string    `json:"status"`
	Accrual    uint      `json:"accrual"`
	UploadedAt time.Time `json:"uploaded_at"`
	User       string    `json:"user_id"`
}
