package helpers

import (
	"time"
)

func NowPlusSeconds(seconds int) time.Time {
	return time.Now().Local().Add(time.Second * time.Duration(seconds))
}

func ParseDateIso(date string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05-0700", date)
}

func FormatDateIso(date time.Time) string {
	return date.Format("2006-01-02T15:04:05-0700")
}
