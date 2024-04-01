package cronroutine

import "time"

// ctime is an unfortunate abstraction that we need to have because time.Time
// only has a field for unix time.
type ctime struct {
	minute     int
	hour       int
	dayOfMonth int
	month      int
	dayOfWeek  int
	year       int
}

func newCTime(time time.Time) *ctime {
	return &ctime{
		minute:     time.Minute(),
		hour:       time.Hour(),
		dayOfMonth: time.Day(),
		month:      int(time.Month()),
		dayOfWeek:  int(time.Weekday()),
		year:       time.Year(),
	}
}

func (ct *ctime) time() time.Time {
	return time.Date(ct.year, time.Month(ct.month), ct.dayOfMonth, ct.hour, ct.minute, second, nanosecond, time.UTC)
}
