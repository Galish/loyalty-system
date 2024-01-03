package entity

import "time"

const TimeLayout = "2006-01-02T15:04:05-07:00"

type Time time.Time

func (t Time) Format() string {
	return time.Time(t).Format(TimeLayout)
}

func (t Time) Value() time.Time {
	return time.Time(t).Round(time.Microsecond)
}
