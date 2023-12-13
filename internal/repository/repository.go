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
	GetUserBalance(context.Context, string) (*Balance, error)
	UpdateBalance(context.Context, string, float32) error

	CreateWithdraw(context.Context, *Withdraw) error
	GetWithdrawals(context.Context, string) ([]*Withdraw, error)
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
	Accrual    float32   `json:"accrual"`
	UploadedAt time.Time `json:"uploaded_at"`
	User       string    `json:"user_id"`
}

type Balance struct {
	User      string    `json:"user_id"`
	Current   float32   `json:"current"`
	Withdrawn float32   `json:"withdrawn"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Withdraw struct {
	Order       string
	User        string
	Sum         float32
	ProcessedAt time.Time
}
