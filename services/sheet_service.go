package services

import (
	"sync"
	"time"
)

type MemorySheet struct {
	Date time.Time
	Year int
	Text string
}

type SheetService struct {
	dir string
	mu  sync.Mutex
}

func NewSheetService(dir string) *SheetService {
	return &SheetService{
		dir: dir,
	}
}

func (s *SheetService) ReadDir() ([]*MemorySheet, error) {
	panic("//2d")
}

func (s *SheetService) LookUpSheets(sheets []*MemorySheet, date time.Time) ([]*MemorySheet, error) {
	panic("//2d")
}

func (s *SheetService) WriteSheet(date time.Time, content string) error {
	panic("//2d")
}

func (s *SheetService) DeleteSheet(date time.Time) error {
	panic("//2d")
}

type RemindPattern []int

func IsDateReminding(date time.Time, today time.Time, p RemindPattern) bool {
	step := 0
	for {
		distance := p[min(step, len(p))]
		date = date.AddDate(0, 0, distance)
		if date.Equal(today) {
			return true
		}
		if date.After(today) {
			return false
		}
		step++
	}
}
