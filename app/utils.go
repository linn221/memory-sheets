package app

import "time"

func Today() time.Time {
	// local time.Now with time part zero
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}
