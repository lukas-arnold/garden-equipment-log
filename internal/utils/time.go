package utils

import "time"

func ParseDateTime(value string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04", value)
}
