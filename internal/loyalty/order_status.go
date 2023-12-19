package loyalty

type Status string

const (
	StatusNew        = Status("NEW")
	StatusRegistered = Status("REGISTERED")
	StatusProcessing = Status("PROCESSING")
	StatusInvalid    = Status("INVALID")
	StatusProcessed  = Status("PROCESSED")

	TimeLayout = "2006-01-02T15:04:05-07:00"
)

func (s Status) isFinal() bool {
	switch s {
	case StatusInvalid, StatusProcessed:
		return true

	default:
		return false
	}
}
