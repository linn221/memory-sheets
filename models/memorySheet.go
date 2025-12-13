package models

import "time"

type MemorySheet struct {
	Date time.Time
	Year int
	Text string
}

func (s *MemorySheet) DisplayDate() string {
	return s.Date.Format(time.DateOnly)
}
