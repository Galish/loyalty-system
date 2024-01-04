package entity

import "github.com/Galish/loyalty-system/internal/validation"

type Balance struct {
	User      string
	Current   float32
	Withdrawn float32
	UpdatedAt Time
}

type Withdrawal struct {
	Order       string
	Sum         float32
	User        string
	ProcessedAt Time
}

func (w Withdrawal) IsValid() bool {
	return validation.IsValidLuhn(w.Order)
}

type Enrollment struct {
	User        string
	Sum         float32
	ProcessedAt Time
}
