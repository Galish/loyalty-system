package entity

import (
	"github.com/Galish/loyalty-system/internal/app/validation"
)

const (
	StatusNew        = Status("NEW")
	StatusRegistered = Status("REGISTERED")
	StatusProcessing = Status("PROCESSING")
	StatusInvalid    = Status("INVALID")
	StatusProcessed  = Status("PROCESSED")
)

type Order struct {
	ID         string
	Status     Status
	Accrual    float32
	UploadedAt Time
	User       string
}

func (o Order) Validate() error {
	return validation.LuhnValidate(o.ID)
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
