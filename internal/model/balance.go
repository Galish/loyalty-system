package model

import "time"

type Balance struct {
	User      string
	Current   float32
	Withdrawn float32
	UpdatedAt time.Time
}

type Withdrawal struct {
	Order       OrderNumber
	Sum         float32
	User        string
	ProcessedAt time.Time
}

type Enrollment struct {
	User        string
	Sum         float32
	ProcessedAt time.Time
}
