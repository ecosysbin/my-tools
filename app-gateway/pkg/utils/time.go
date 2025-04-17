package utils

import "time"

func TimeNow() *time.Time {
	timeNow := time.Now()
	return &timeNow
}

func FormartTimeNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse("2006-01-02", timeStr)
}

func ParseTimeEvent(event string) string {
	timeNow := FormartTimeNow()
	return timeNow + " " + event
}
