package entity

import (
	"time"

	"github.com/Galish/loyalty-system/internal/validation"
)

type Balance struct {
	User      string
	Current   float32
	Withdrawn float32
	UpdatedAt time.Time
}

type Withdrawal struct {
	Order       string
	Sum         float32
	User        string
	ProcessedAt time.Time
}

func (w Withdrawal) IsValid() bool {
	return validation.IsValidLuhn(w.Order)
}

type Enrollment struct {
	User        string
	Sum         float32
	ProcessedAt time.Time
}
