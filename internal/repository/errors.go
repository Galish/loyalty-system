package repository

import "errors"

var (
	ErrOrderExists       = errors.New("order has already been added")
	ErrOrderConflict     = errors.New("order has already been added by another user")
	ErrUserConflict      = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrInsufficientFunds = errors.New("insufficient funds in the account")
)
