package entity

import (
	"github.com/Galish/loyalty-system/internal/app/validation"
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
