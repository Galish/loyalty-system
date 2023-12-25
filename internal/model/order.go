package model

import (
	"github.com/ShiraazMoollatjie/goluhn"
)

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
