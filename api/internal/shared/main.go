package shared

import "time"

type LocalNowTimeProvider = func() time.Time
type UTCNowTimeProvider = func() time.Time

func GetTimeZoneIndependantValue(datetime time.Time, tz string) time.Time {
	// change TZ without changing actual time value
	loc, _ := time.LoadLocation(tz)
	return time.Date(
		datetime.Year(),
		datetime.Month(),
		datetime.Day(),
		datetime.Hour(),
		datetime.Minute(),
		datetime.Second(),
		datetime.Nanosecond(),
		loc,
	)
}
