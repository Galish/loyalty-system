package model

import "time"

type Time time.Time

func (t Time) Format() string {
	return time.Time(t).Format(TimeLayout)
}

func (t Time) Value() time.Time {
	return time.Time(t).Round(time.Microsecond)
}
