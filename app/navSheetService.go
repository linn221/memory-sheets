package app

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/linn221/memory-sheets/models"
)

type NavSheetService struct {
	mu     sync.Mutex
	dir    string
	sheets []*models.NavSheet
}

// ReadDir reads the nav directory and scans markdown files, storing them in NavSheetService
// The filename without extension becomes the Title of NavSheet
func (s *NavSheetService) ReadDir() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sheets = []*models.NavSheet{}

	err := filepath.Walk(s.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only process .md files
		if !info.IsDir() && strings.HasSuffix(path, ".md") {
			sheet, err := parseFilepathToNavSheet(s.dir, path)
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

// Create creates a new NavSheet with the given title and text
// Writes to file in nav directory and updates in-memory sheets
func (s *NavSheetService) Create(title string, text string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	filePath := s.fromTitleToFilepath(title)

	// Check if file exists
	if fileExists(filePath) {
		return fmt.Errorf("nav sheet already exists with title %s", title)
	}

	// Ensure directory exists (including subdirectories if title contains path separators)
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Write the file
	if err := writeFileContent(filePath, text); err != nil {
		return err
	}

	// Update in-memory sheet
	for _, sheet := range s.sheets {
		if sheet.Title == title {
			sheet.Text = text
			return nil
		}
	}

	// If not found in memory, add it
	sheet := &models.NavSheet{
		Title: title,
		Text:  text,
	}
	s.sheets = append(s.sheets, sheet)

	return nil
}

// Update updates an existing NavSheet with the given title and text
// Writes to file in nav directory and updates in-memory sheets
func (s *NavSheetService) Update(title string, text string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	filePath := s.fromTitleToFilepath(title)

	// Check if file exists
	if !fileExists(filePath) {
		return fmt.Errorf("nav sheet does not exist with title %s", title)
	}

	// Write the file
	if err := writeFileContent(filePath, text); err != nil {
		return err
	}

	// Update in-memory sheet
	for _, sheet := range s.sheets {
		if sheet.Title == title {
			sheet.Text = text
			return nil
		}
	}

	// If not found in memory, add it
	sheet := &models.NavSheet{
		Title: title,
		Text:  text,
	}
	s.sheets = append(s.sheets, sheet)

	return nil
}

// Delete deletes a NavSheet with the given title
// Deletes the file and removes from in-memory sheets
func (s *NavSheetService) Delete(title string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	filePath := s.fromTitleToFilepath(title)

	// Check if file exists
	if !fileExists(filePath) {
		return fmt.Errorf("nav sheet does not exist with title %s", title)
	}

	// Delete the file
	if err := deleteFile(filePath); err != nil {
		return err
	}

	// Remove from in-memory sheets
	for i, sheet := range s.sheets {
		if sheet.Title == title {
			s.sheets = append(s.sheets[:i], s.sheets[i+1:]...)
			return nil
		}
	}

	return nil
}

// Get retrieves a NavSheet by title
// Returns the sheet from memory or loads from file if not in memory
func (s *NavSheetService) Get(title string) (*models.NavSheet, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check in-memory sheets first
	for _, sheet := range s.sheets {
		if sheet.Title == title {
			return sheet, nil
		}
	}

	// If not in memory, try to load from file
	filePath := s.fromTitleToFilepath(title)
	if fileExists(filePath) {
		sheet, err := parseFilepathToNavSheet(s.dir, filePath)
		if err != nil {
			return nil, err
		}
		// Add to in-memory sheets
		s.sheets = append(s.sheets, sheet)
		return sheet, nil
	}

	return nil, fmt.Errorf("nav sheet does not exist with title %s", title)
}

// ListSheets returns all NavSheets
func (s *NavSheetService) ListSheets() []*models.NavSheet {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Return a copy of the slice
	sheets := make([]*models.NavSheet, len(s.sheets))
	copy(sheets, s.sheets)
	return sheets
}

// Search searches through all nav sheets using the provided regex pattern
// Returns matching sheets with Text field modified to bold matched strings using markdown (**text**)
func (s *NavSheetService) Search(patternStr string) ([]*models.NavSheet, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if patternStr == "" {
		return s.sheets, nil
	}

	// Compile regex pattern with case-insensitive flag
	re, err := regexp.Compile("(?i)" + patternStr)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %v", err)
	}

	var matchingSheets []*models.NavSheet
	for _, sheet := range s.sheets {
		// Check if the pattern matches the text or title
		if re.MatchString(sheet.Text) || re.MatchString(sheet.Title) {
			// Create a copy of the sheet with highlighted matches
			highlightedText := re.ReplaceAllStringFunc(sheet.Text, func(match string) string {
				return "**" + match + "**"
			})
			highlightedTitle := re.ReplaceAllStringFunc(sheet.Title, func(match string) string {
				return "**" + match + "**"
			})

			// Create a new sheet with highlighted text
			highlightedSheet := &models.NavSheet{
				Title: highlightedTitle,
				Text:  highlightedText,
			}
			matchingSheets = append(matchingSheets, highlightedSheet)
		}
	}

	return matchingSheets, nil
}

// fromTitleToFilepath converts a title to a filepath
// The title becomes the filename with .md extension in the nav directory
func (s *NavSheetService) fromTitleToFilepath(title string) string {
	filename := title + ".md"
	return filepath.Join(s.dir, filename)
}

// parseFilepathToNavSheet parses a filepath into a NavSheet
// The filename without extension becomes the Title, and the file content becomes the Text
func parseFilepathToNavSheet(dir string, path string) (*models.NavSheet, error) {
	// Get relative path from base directory
	relPath, err := filepath.Rel(dir, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get relative path: %v", err)
	}

	// Remove .md extension to get title
	title := strings.TrimSuffix(relPath, ".md")

	// Read file content
	content, err := readFileContent(path)
	if err != nil {
		return nil, err
	}

	return &models.NavSheet{
		Title: title,
		Text:  content,
	}, nil
}
