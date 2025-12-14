package models

import "time"

type MemorySheet struct {
	Date time.Time
	Year int
	Text string
}

func (s *MemorySheet) Url() string {
	return "sheets/" + s.Date.Format(time.DateOnly)
}

func (s *MemorySheet) DateStr() string {
	return s.Date.Format(time.DateOnly)

}

func (s *MemorySheet) Title() string {
	return s.Date.Format("Jan 2")
}
