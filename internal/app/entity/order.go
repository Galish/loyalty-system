package entity

import (
	"errors"

	"github.com/ShiraazMoollatjie/goluhn"
)

var ErrInvalidOrderNumber = errors.New("invalid order number value")

const (
	StatusNew        = Status("NEW")
	StatusRegistered = Status("REGISTERED")
	StatusProcessing = Status("PROCESSING")
	StatusInvalid    = Status("INVALID")
	StatusProcessed  = Status("PROCESSED")
)

type Order struct {
	ID         OrderNumber
	Status     Status
	Accrual    float32
	UploadedAt Time
	User       string
}

func (o Order) Validate() error {
	if !o.ID.IsValid() {
		return ErrInvalidOrderNumber
	}

	return nil
}

type Status string

func (s Status) IsFinal() bool {
	switch s {
	case StatusInvalid, StatusProcessed:
		return true

	default:
		return false
	}
}

type OrderNumber string

func (num OrderNumber) IsValid() bool {
	if num.String() == "" {
		return false
	}

	if err := goluhn.Validate(num.String()); err != nil {
		return false
	}

	return true
}

func (num OrderNumber) String() string {
	return string(num)
}
