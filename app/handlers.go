package app

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/linn221/memory-sheets/models"
	"github.com/linn221/memory-sheets/views"
)

// ShowTodaySheets handles GET /sheets - returns sheets that need to be reminded today
func (a *App) ShowTodaySheets(vr *views.ViewRenderer) error {
	today := Today()
	remindingSheets, err := a.sheetService.LookUpSheets(today)
	if err != nil {
		return err
	}
	todaySheet, err := a.sheetService.GetSheetByDate(today)
	if err != nil {
		todaySheet = nil
	}

	return vr.IndexPage(remindingSheets, todaySheet, a.navSheetService.ListSheets())
}

// ShowAllSheets handles GET /all-sheets - returns all sheets
func (a *App) ShowAllSheets(vr *views.ViewRenderer) error {
	todaySheet, err := a.sheetService.GetSheetByDate(Today())
	if err != nil {
		todaySheet = nil
	}
	return vr.IndexPage(a.sheetService.sheets, todaySheet, a.navSheetService.ListSheets())
}

// ShowEditSheet handles GET /sheets/{date}/edit - returns the edit page for a sheet
func (a *App) ShowEditSheet(vr *views.ViewRenderer) error {
	r := vr.Request()
	dateStr := r.PathValue("date")
	date, err := time.Parse(time.DateOnly, dateStr)
	if err != nil {
		return err
	}
	sheet, err := a.sheetService.GetSheetByDate(date)
	if err != nil {
		return err
	}
	content := sheet.Text
	return vr.ShowEditSheet(dateStr, content)
}

// ShowSheet handles GET /sheets/{date} - returns a specific sheet by date
func (a *App) ShowSheet(vr *views.ViewRenderer) error {
	r := vr.Request()
	dateStr := r.PathValue("date")
	if dateStr == "" {
		return errors.New("date cannot be empty")
	}
	date, err := time.Parse(time.DateOnly, dateStr)
	if err != nil {
		return err
	}
	today := Today()
	remindingSheets, err := a.sheetService.LookUpSheets(today)
	if err != nil {
		return err
	}
	var current *models.MemorySheet
	for _, rSheet := range remindingSheets {
		if rSheet.Date.Equal(date) {
			current = rSheet
			break
		}
	}
	if current == nil {
		return errors.New("note not found")
	}
	return vr.SheetComponent(current)
}

func (a *App) HandleCreateSheet(vr *views.ViewRenderer) error {
	r := vr.Request()

	// Read request body
	content := r.FormValue("content")

	// Update the sheet
	err := a.sheetService.CreateSheet(content)
	if err != nil {
		return err
	}

	sheet, err := a.sheetService.GetSheetByDate(Today())
	if err != nil {
		return err
	}
	return vr.SheetComponent(sheet)
}

// HandleUpdateSheet handles PUT /sheets/{date} - updates an existing sheet
func (a *App) HandleUpdateSheet(vr *views.ViewRenderer) error {
	r := vr.Request()
	dateStr := r.PathValue("date")
	if dateStr == "" {
		return errors.New("date cannot be empty")
	}

	date, err := time.Parse(time.DateOnly, dateStr)
	if err != nil {
		return errors.New("invalid date format")
	}

	// Read request body
	content := r.FormValue("content")

	// Update the sheet
	err = a.sheetService.UpdateSheet(date, content)
	if err != nil {
		return err
	}

	sheet, err := a.sheetService.GetSheetByDate(date)
	if err != nil {
		return err
	}
	return vr.SheetComponent(sheet)
}

// HandleDeleteSheet handles DELETE /sheets/{date} - deletes a sheet
func (a *App) HandleDeleteSheet(vr *views.ViewRenderer) error {
	r := vr.Request()
	dateStr := r.PathValue("date")
	if dateStr == "" {
		return errors.New("date cannot be empty")
	}

	date, err := time.Parse(time.DateOnly, dateStr)
	if err != nil {
		return errors.New("invalid date format")
	}

	// Delete the sheet
	err = a.sheetService.DeleteSheet(date)
	if err != nil {
		return err
	}

	return nil
}

// ShowChangePattern handles GET /change-pattern - shows the pattern editor
func (a *App) ShowChangePattern(vr *views.ViewRenderer) error {
	pattern := a.sheetService.GetPattern()
	// Convert pattern to a map for easy lookup
	selectedMap := make(map[int]bool)
	for _, day := range pattern {
		if day >= 1 && day <= 200 {
			selectedMap[day] = true
		}
	}
	return vr.ShowChangePattern(selectedMap)
}

// HandlePostChangePattern handles POST /change-pattern - saves the pattern
func (a *App) HandlePostChangePattern(vr *views.ViewRenderer) error {
	r := vr.Request()
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("failed to parse form: %v", err)
	}

	// Collect selected days from form
	selectedDays := make(map[int]bool)
	for key := range r.PostForm {
		if strings.HasPrefix(key, "day_") {
			var day int
			if _, err := fmt.Sscanf(key, "day_%d", &day); err == nil {
				if day >= 1 && day <= 200 {
					selectedDays[day] = true
				}
			}
		}
	}

	// Convert to sorted slice
	pattern := make(RemindPattern, 0, len(selectedDays))
	for day := 1; day <= 200; day++ {
		if selectedDays[day] {
			pattern = append(pattern, day)
		}
	}

	// Save to JSON file
	patternFile := "pattern.json"
	if err := SavePatternToJSON(patternFile, pattern); err != nil {
		return fmt.Errorf("failed to save pattern: %v", err)
	}

	// Update in-memory pattern
	a.sheetService.UpdatePattern(pattern)

	// Redirect to /change-pattern
	http.Redirect(vr.ResponseWriter(), vr.Request(), "/change-pattern", http.StatusSeeOther)
	return nil
}

// HandleSearch handles GET /search - searches sheets using regex pattern from q query parameter
func (a *App) HandleSearch(vr *views.ViewRenderer) error {
	r := vr.Request()
	query := r.URL.Query().Get("q")

	// Search for matching memory sheets
	memoryResults, err := a.sheetService.Search(query)
	if err != nil {
		return err
	}

	// Search for matching nav sheets
	navResults, err := a.navSheetService.Search(query)
	if err != nil {
		return err
	}

	return vr.SearchResults(memoryResults, navResults)
}

// ShowNavSheet handles GET /nav-sheets/{title} - returns a specific nav sheet
func (a *App) ShowNavSheet(vr *views.ViewRenderer) error {
	r := vr.Request()
	title := r.PathValue("title")
	if title == "" {
		return errors.New("title cannot be empty")
	}

	sheet, err := a.navSheetService.Get(title)
	if err != nil {
		return err
	}
	return vr.NavSheetComponent(sheet)
}

// ShowEditNavSheet handles GET /nav-sheets/{title}/edit - returns the edit page for a nav sheet
func (a *App) ShowEditNavSheet(vr *views.ViewRenderer) error {
	r := vr.Request()
	title := r.PathValue("title")
	if title == "" {
		return errors.New("title cannot be empty")
	}

	sheet, err := a.navSheetService.Get(title)
	if err != nil {
		return err
	}
	content := sheet.Text
	return vr.ShowEditNavSheet(title, content)
}

func (a *App) ShowCreateNavSheet(vr *views.ViewRenderer) error {
	return vr.ShowCreateNavSheet()
}

// HandleUpdateNavSheet handles PUT /nav-sheets/{title} - updates an existing nav sheet
func (a *App) HandleUpdateNavSheet(vr *views.ViewRenderer) error {
	r := vr.Request()
	title := r.PathValue("title")
	if title == "" {
		return errors.New("title cannot be empty")
	}

	// Read request body
	content := r.FormValue("content")

	// Update the sheet
	err := a.navSheetService.Update(title, content)
	if err != nil {
		return err
	}

	sheet, err := a.navSheetService.Get(title)
	if err != nil {
		return err
	}
	return vr.NavSheetComponent(sheet)
}

// HandleDeleteNavSheet handles DELETE /nav-sheets/{title} - deletes a nav sheet
func (a *App) HandleDeleteNavSheet(vr *views.ViewRenderer) error {
	r := vr.Request()
	title := r.PathValue("title")
	if title == "" {
		return errors.New("title cannot be empty")
	}

	// Delete the sheet
	err := a.navSheetService.Delete(title)
	if err != nil {
		return err
	}

	return vr.NavSheetsComponent(a.navSheetService.ListSheets())
}

func (a *App) HandleCreateNavSheet(vr *views.ViewRenderer) error {
	r := vr.Request()
	title := r.FormValue("title")
	if title == "" {
		return errors.New("title cannot be empty")
	}
	content := r.FormValue("content")

	// Delete the sheet
	err := a.navSheetService.Create(title, content)
	if err != nil {
		return err
	}

	sheet, err := a.navSheetService.Get(title)
	if err != nil {
		return err
	}
	if err := vr.NavSheetComponent(sheet); err != nil {
		return err
	}
	return vr.NavSheetsComponent(a.navSheetService.ListSheets())
}
