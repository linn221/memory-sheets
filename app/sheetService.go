package app

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/linn221/memory-sheets/models"
)

type SheetService struct {
	mu      sync.Mutex
	pattern RemindPattern
	dir     string
	sheets  []*models.MemorySheet
}

// read the dir directory and scan sheets []*models.MemorySheet, store the sheets in SheetService
// 2025/jan-1.md will be turned into models.MemorySheet of date (jan 1 2025 00:00:00 UTC), and Text will be the contents of the file
// the sheets slice will always be ordered by Date field ascending
func (s *SheetService) ReadDir() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sheets = []*models.MemorySheet{}

	err := filepath.Walk(s.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only process .md files
		if !info.IsDir() && strings.HasSuffix(path, ".md") {
			sheet, err := parseFilepathToSheet(s.dir, path)
			if err != nil {
				fmt.Printf("%s file does not get parsed for some reason: %v\n", path, err)
			} else {
				s.sheets = append(s.sheets, sheet)
			}
		}

		return nil
	})

	if err == nil {
		s.sortSheets()
	}

	return err
}

func (s *SheetService) LookUpSheets(date time.Time) ([]*models.MemorySheet, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var remindingSheets []*models.MemorySheet
	for _, sheet := range s.sheets {
		sheetDate := sheet.Date
		if IsDateReminding(sheetDate, date, s.pattern) {
			remindingSheets = append(remindingSheets, sheet)
		}
	}
	return remindingSheets, nil
}

// will create new file if it does not exists with the name monthPrefix-day.md
func (s *SheetService) CreateSheet(date time.Time, content string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := s.fromDateToFilepath(date)

	// Check if file already exists
	if fileExists(path) {
		return fmt.Errorf("sheet already exists for date %s", date.Format(time.DateOnly))
	}

	// Write the file
	if err := writeFileContent(path, content); err != nil {
		return err
	}

	// Add to in-memory sheets
	normalizedDate := normalizeDate(date)
	sheet := &models.MemorySheet{
		Date: normalizedDate,
		Year: date.Year(),
		Text: content,
	}
	s.insertSheetInOrder(sheet)

	return nil
}

// update text file if it exists
func (s *SheetService) UpdateSheet(date time.Time, content string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	filepath := s.fromDateToFilepath(date)

	// Check if file exists
	if !fileExists(filepath) {
		return fmt.Errorf("sheet does not exist for date %s", date.Format(time.DateOnly))
	}

	// Write the file
	if err := writeFileContent(filepath, content); err != nil {
		return err
	}

	// Update in-memory sheet
	normalizedDate := normalizeDate(date)
	for _, sheet := range s.sheets {
		if sheet.Date.Equal(normalizedDate) {
			sheet.Text = content
			return nil
		}
	}

	// If not found in memory, add it
	sheet := &models.MemorySheet{
		Date: normalizedDate,
		Year: date.Year(),
		Text: content,
	}
	s.insertSheetInOrder(sheet)

	return nil
}

// delete the file
func (s *SheetService) DeleteSheet(date time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	filepath := s.fromDateToFilepath(date)

	// Check if file exists
	if !fileExists(filepath) {
		return fmt.Errorf("sheet does not exist for date %s", date.Format(time.DateOnly))
	}

	// Delete the file
	if err := deleteFile(filepath); err != nil {
		return err
	}

	// Remove from in-memory sheets
	normalizedDate := normalizeDate(date)
	for i, sheet := range s.sheets {
		if sheet.Date.Equal(normalizedDate) {
			s.sheets = append(s.sheets[:i], s.sheets[i+1:]...)
			return nil
		}
	}

	return nil
}

// read the sheet
// returns the index of the sheet in Sheets slice
func (s *SheetService) GetSheetByDate(date time.Time) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check in-memory sheets first
	normalizedDate := normalizeDate(date)
	for i, sheet := range s.sheets {
		if sheet.Date.Equal(normalizedDate) {
			return i, nil
		}
	}

	// If not in memory, try to load from file
	filepath := s.fromDateToFilepath(date)
	if fileExists(filepath) {
		index, err := s.loadSheetFromFile(date, filepath)
		if err != nil {
			return -1, err
		}
		return index, nil
	}

	return -1, fmt.Errorf("sheet does not exist for date %s", date.Format(time.DateOnly))
}

// read the sheet
func (s *SheetService) IsSheetExist(date time.Time) bool {
	filepath := s.fromDateToFilepath(date)

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check in-memory sheets first
	normalizedDate := normalizeDate(date)
	for _, sheet := range s.sheets {
		if sheet.Date.Equal(normalizedDate) {
			return true
		}
	}

	// If not in memory, check file system
	return fileExists(filepath)
}

// convert date to filepath, jan 01 2025 of time will be jan-01.md in the service directory
// make use of this method in Create/Update/DeleteSheet methods
func (s *SheetService) fromDateToFilepath(date time.Time) string {
	month := strings.ToLower(date.Format("Jan"))
	day := date.Day()
	filename := fmt.Sprintf("%s-%02d.md", month, day)
	return filepath.Join(s.dir, filename)
}

// sortSheets sorts the sheets slice by Date field ascending
func (s *SheetService) sortSheets() {
	sort.Slice(s.sheets, func(i, j int) bool {
		return s.sheets[i].Date.Before(s.sheets[j].Date)
	})
}

// insertSheetInOrder inserts a sheet into the sheets slice at the correct position to maintain ascending order by Date
func (s *SheetService) insertSheetInOrder(sheet *models.MemorySheet) {
	// Find the insertion point
	insertIndex := sort.Search(len(s.sheets), func(i int) bool {
		return !s.sheets[i].Date.Before(sheet.Date)
	})

	// Insert at the found position
	if insertIndex == len(s.sheets) {
		s.sheets = append(s.sheets, sheet)
	} else {
		// Extend slice by one
		s.sheets = append(s.sheets, nil)
		// Shift elements to the right
		copy(s.sheets[insertIndex+1:], s.sheets[insertIndex:])
		// Insert the new sheet
		s.sheets[insertIndex] = sheet
	}
}

// loadSheetFromFile reads a sheet from the file system and stores it in the slice at its correct position
// Returns the index where the sheet was inserted
func (s *SheetService) loadSheetFromFile(date time.Time, filepath string) (int, error) {
	content, err := readFileContent(filepath)
	if err != nil {
		return -1, err
	}

	normalizedDate := normalizeDate(date)
	sheet := &models.MemorySheet{
		Date: normalizedDate,
		Year: date.Year(),
		Text: content,
	}

	// Find the insertion point before inserting
	insertIndex := sort.Search(len(s.sheets), func(i int) bool {
		return !s.sheets[i].Date.Before(normalizedDate)
	})

	// Insert at the found position
	if insertIndex == len(s.sheets) {
		s.sheets = append(s.sheets, sheet)
	} else {
		// Extend slice by one
		s.sheets = append(s.sheets, nil)
		// Shift elements to the right
		copy(s.sheets[insertIndex+1:], s.sheets[insertIndex:])
		// Insert the new sheet
		s.sheets[insertIndex] = sheet
	}

	return insertIndex, nil
}

type RemindPattern []int

func IsDateReminding(date time.Time, today time.Time, p RemindPattern) bool {
	step := 0
	for {
		distance := p[min(step, len(p)-1)]
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
