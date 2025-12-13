package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/linn221/memory-sheets/models"
)

type SheetService struct {
	mu      sync.Mutex
	pattern RemindPattern
	sheets  []*models.MemorySheet
}

// read the dir directory and scan sheets []*models.MemorySheet, store the sheets in SheetService
// 2025/jan-1.md will be turned into models.MemorySheet of date (jan 1 2025 00:00:00 UTC), and Text will be the contents of the file
func (s *SheetService) ReadDir(dir string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sheets = []*models.MemorySheet{}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only process .md files
		if !info.IsDir() && strings.HasSuffix(path, ".md") {
			sheet, err := parseFilepathToSheet(dir, path)
			if err != nil {
				fmt.Printf("%s file does not get parsed for some reason: %v\n", path, err)
			} else {
				s.sheets = append(s.sheets, sheet)
			}
		}

		return nil
	})

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

	path := fromDateToFilepath(date)

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
	s.sheets = append(s.sheets, sheet)

	return nil
}

// update text file if it exists
func (s *SheetService) UpdateSheet(date time.Time, content string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	filepath := fromDateToFilepath(date)

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
	s.sheets = append(s.sheets, sheet)

	return nil
}

// delete the file
func (s *SheetService) DeleteSheet(date time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	filepath := fromDateToFilepath(date)

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
func (s *SheetService) GetSheetByDate(date time.Time) (*models.MemorySheet, error) {
	filepath := fromDateToFilepath(date)

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check in-memory sheets first
	normalizedDate := normalizeDate(date)
	for _, sheet := range s.sheets {
		if sheet.Date.Equal(normalizedDate) {
			return sheet, nil
		}
	}

	// If not in memory, try to read from file
	if fileExists(filepath) {
		content, err := readFileContent(filepath)
		if err != nil {
			return nil, err
		}

		sheet := &models.MemorySheet{
			Date: normalizedDate,
			Year: date.Year(),
			Text: content,
		}

		// Add to in-memory sheets
		s.sheets = append(s.sheets, sheet)

		return sheet, nil
	}

	return nil, fmt.Errorf("sheet does not exist for date %s", date.Format(time.DateOnly))
}

// read the sheet
func (s *SheetService) IsSheetExist(date time.Time) bool {
	filepath := fromDateToFilepath(date)

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

// normalizeDate normalizes a date to midnight UTC (00:00:00 UTC)
func normalizeDate(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
}

// convert date to filepath, jan 01 2025 of time will be 2025/jan-01.md
// make use of this function in Create/Update/DeleteSheet methods
func fromDateToFilepath(date time.Time) string {
	year := date.Year()
	month := strings.ToLower(date.Format("Jan"))
	day := date.Day()
	return fmt.Sprintf("%d/%s-%02d.md", year, month, day)
}

// readFileContent reads the content of a file at the given path
func readFileContent(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// writeFileContent writes content to a file at the given path, creating parent directories if needed
func writeFileContent(path string, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

// deleteFile deletes a file at the given path
func deleteFile(path string) error {
	return os.Remove(path)
}

// fileExists checks if a file exists at the given path
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// parseFilepathToSheet parses a filepath like "sheets/2025/jan-01.md" into a MemorySheet
func parseFilepathToSheet(dir string, path string) (*models.MemorySheet, error) {
	// Remove the base directory and .md extension
	relPath := strings.TrimPrefix(path, dir+"/")
	relPath = strings.TrimSuffix(relPath, ".md")

	// Parse format: YYYY/month-day
	parts := strings.Split(relPath, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid filepath format: %s", path)
	}

	var year int
	if _, err := fmt.Sscanf(parts[0], "%d", &year); err != nil {
		return nil, fmt.Errorf("invalid year in filepath: %s", path)
	}

	// Parse month-day format
	var monthStr string
	var day int
	if _, err := fmt.Sscanf(parts[1], "%3s-%d", &monthStr, &day); err != nil {
		return nil, fmt.Errorf("invalid month-day format in filepath: %s", path)
	}

	// Convert month abbreviation to month number
	monthMap := map[string]time.Month{
		"jan": time.January, "feb": time.February, "mar": time.March,
		"apr": time.April, "may": time.May, "jun": time.June,
		"jul": time.July, "aug": time.August, "sep": time.September,
		"oct": time.October, "nov": time.November, "dec": time.December,
	}

	month, ok := monthMap[strings.ToLower(monthStr)]
	if !ok {
		return nil, fmt.Errorf("invalid month abbreviation: %s", monthStr)
	}

	date := normalizeDate(time.Date(year, month, day, 0, 0, 0, 0, time.UTC))

	// Read file content
	content, err := readFileContent(path)
	if err != nil {
		return nil, err
	}

	return &models.MemorySheet{
		Date: date,
		Year: year,
		Text: content,
	}, nil
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
