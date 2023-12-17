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

type OrderRepository interface {
	CreateOrder(context.Context, *Order) error
	UserOrders(context.Context, string) ([]*Order, error)
	UpdateOrder(context.Context, *Order) error
}

type BalanceRepository interface {
	UserBalance(context.Context, string) (*Balance, error)
	Enroll(context.Context, *Enrollment) error
	Withdraw(context.Context, *Withdrawal) error
	Withdrawals(context.Context, string) ([]*Withdrawal, error)
}

type LoyaltyRepository interface {
	OrderRepository
	BalanceRepository
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

type Enrollment struct {
	User        string
	Sum         float32
	ProcessedAt time.Time
}

type Withdrawal struct {
	Order       string
	User        string
	Sum         float32
	ProcessedAt time.Time
}
