package entity

const (
	StatusNew        = Status("NEW")
	StatusRegistered = Status("REGISTERED")
	StatusProcessing = Status("PROCESSING")
	StatusInvalid    = Status("INVALID")
	StatusProcessed  = Status("PROCESSED")
)

type Status string

func (s Status) IsFinal() bool {
	switch s {
	case StatusInvalid, StatusProcessed:
		return true

	default:
		return false
	}
}
