package helpers

import "time"

func StartOfDay(dateTime *time.Time) *time.Time {
	output := time.Date(
		dateTime.Year(),
		dateTime.Month(),
		dateTime.Day(),
		0, 0, 0, 0,
		dateTime.Location(),
	)

	return &output
}

func EndOfDay(dateTime *time.Time) *time.Time {
	output := time.Date(
		dateTime.Year(),
		dateTime.Month(),
		dateTime.Day(),
		23, 59, 59, 999999999,
		dateTime.Location(),
	)

	return &output
}

func TimeZoneIndependentTimeNow(locationTimeZone string) time.Time {
	timeZone, _ := time.LoadLocation(locationTimeZone)
	utcTimeZone, _ := time.LoadLocation("UTC")
	now := time.Now().In(timeZone)
	return time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second(),
		now.Nanosecond(),
		utcTimeZone)
}
