package repository

import "errors"

var (
	ErrInsufficientFunds = errors.New("insufficient funds in the account")
	ErrNothingFound      = errors.New("nothing was found")
	ErrOrderConflict     = errors.New("order has already been added by another user")
	ErrOrderExists       = errors.New("order has already been added")
	ErrUserConflict      = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)
