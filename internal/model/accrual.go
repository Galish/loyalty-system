package model

type Accrual struct {
	Order  OrderNumber
	Status Status
	Value  float32
	User   string
}
