package entity

type Balance struct {
	User      string
	Current   float32
	Withdrawn float32
	UpdatedAt Time
}

type Withdrawal struct {
	Order       OrderNumber
	Sum         float32
	User        string
	ProcessedAt Time
}

type Enrollment struct {
	User        string
	Sum         float32
	ProcessedAt Time
}
