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
