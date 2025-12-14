package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/linn221/memory-sheets/models"
)

func Today() time.Time {
	// local time.Now with time part zero
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}

// normalizeDate normalizes a date to midnight UTC (00:00:00 UTC)
func normalizeDate(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
}

// readFileContent reads the content of a file at the given path
func readFileContent(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// writeFileContent writes content to a file at the given path
func writeFileContent(path string, content string) error {
	// dir := filepath.Dir(path)
	// if err := os.MkdirAll(dir, 0755); err != nil {
	// 	return err
	// }
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

// parseFilepathToSheet parses a filepath like "sheets/jan-01.md" or "sheets/2025/jan-01.md" into a MemorySheet
func parseFilepathToSheet(dir string, path string) (*models.MemorySheet, error) {
	// Get relative path from base directory
	relPath, err := filepath.Rel(dir, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get relative path: %v", err)
	}
	// Remove .md extension
	relPath = strings.TrimSuffix(relPath, ".md")

	var year int
	var monthStr string
	var day int

	// Try to parse as format: YYYY/month-day (old format with subdirectory)
	parts := strings.Split(relPath, string(filepath.Separator))
	if len(parts) == 2 {
		// Old format: YYYY/month-day
		if _, err := fmt.Sscanf(parts[0], "%d", &year); err != nil {
			return nil, fmt.Errorf("invalid year in filepath: %s", path)
		}
		if _, err := fmt.Sscanf(parts[1], "%3s-%02d", &monthStr, &day); err != nil {
			return nil, fmt.Errorf("invalid month-day format in filepath: %s", path)
		}
	} else if len(parts) == 1 {
		// New format: month-day (directly in dir)
		if _, err := fmt.Sscanf(parts[0], "%3s-%02d", &monthStr, &day); err != nil {
			return nil, fmt.Errorf("invalid month-day format in filepath: %s", path)
		}
		// Get year from file modification time
		info, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("failed to get file info: %v", err)
		}
		year = info.ModTime().Year()
	} else {
		return nil, fmt.Errorf("invalid filepath format: %s", path)
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
