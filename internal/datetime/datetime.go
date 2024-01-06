package datetime

import "time"

const layout = "2006-01-02T15:04:05-07:00"

func Format(t time.Time) string {
	return time.Time(t).Format(layout)
}

func Round(t time.Time) time.Time {
	return time.Time(t).Round(time.Microsecond)
}
