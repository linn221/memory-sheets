package app

import (
	"sync"
	"time"
)

type MemorySheet struct {
	Date time.Time
	Year int
	Text string
}

func (s *MemorySheet) DisplayDate() string {
	return s.Date.Format(time.DateOnly)
}

type SheetService struct {
	mu sync.Mutex
}

func (s *SheetService) ReadDir(dir string) ([]*MemorySheet, error) {
	panic("//2d")
}

func (s *SheetService) LookUpSheets(sheets []*MemorySheet, date time.Time, pattern RemindPattern) ([]*MemorySheet, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var remindingSheets []*MemorySheet
	for _, sheet := range sheets {
		sheetDate := sheet.Date
		if IsDateReminding(sheetDate, date, pattern) {
			remindingSheets = append(remindingSheets, sheet)
		}
	}
	return remindingSheets, nil
}

func (s *SheetService) WriteSheet(date time.Time, content string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	panic("//2d")
}

func (s *SheetService) DeleteSheet(date time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
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
