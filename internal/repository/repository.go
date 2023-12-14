package repository

import (
	"context"
	"errors"
	"time"
)

var (
	ErrOrderExists       = errors.New("order has already been added")
	ErrOrderConflict     = errors.New("order has already been added by another user")
	ErrUserConflict      = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrInsufficientFunds = errors.New("insufficient funds in the account")
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
	UpdateBalance(context.Context, *BalanceEnrollment) error

	Withdraw(context.Context, *Withdraw) error
	GetWithdrawals(context.Context, string) ([]*Withdraw, error)
}

type User struct {
	ID       string
	Login    string
	Password string
	IsActive bool
}

type Order struct {
	ID         string
	Status     string
	Accrual    float32
	UploadedAt time.Time
	User       string
}

type Balance struct {
	User      string
	Current   float32
	Withdrawn float32
	UpdatedAt time.Time
}

type BalanceEnrollment struct {
	User        string
	Sum         float32
	ProcessedAt time.Time
}

type Withdraw struct {
	Order       string
	User        string
	Sum         float32
	ProcessedAt time.Time
}
