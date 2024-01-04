package entity

import (
	"github.com/Galish/loyalty-system/internal/validation"
)

type Order struct {
	ID         string
	Status     Status
	Accrual    float32
	UploadedAt Time
	User       string
}

func (o Order) IsValid() bool {
	return validation.IsValidLuhn(o.ID)
}
