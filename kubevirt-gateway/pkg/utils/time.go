package utils

import "time"

func TimeNow() *time.Time {
	timeNow := time.Now()
	return &timeNow
}
