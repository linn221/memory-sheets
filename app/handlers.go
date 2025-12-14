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

// HandleGetSheets handles GET /sheets - returns sheets that need to be reminded today
func (a *App) HandleGetSheets(vr *views.ViewRenderer) error {
	today := Today()
	remindingSheets, err := a.sheetService.LookUpSheets(today)
	if err != nil {
		return err
	}
	todaySheet, err := a.sheetService.GetSheetByDate(today)
	if err != nil {
		todaySheet = nil
	}

	return vr.IndexPage(&TheSession, remindingSheets, todaySheet)
}

// HandleGetAllSheets handles GET /all-sheets - returns all sheets
func (a *App) HandleGetAllSheets(vr *views.ViewRenderer) error {
	todaySheet, err := a.sheetService.GetSheetByDate(Today())
	if err != nil {
		todaySheet = nil
	}
	return vr.IndexPage(&TheSession, a.sheetService.sheets, todaySheet)
}

// HandleEditSheet handles GET /sheets/{date}/edit - returns the edit page for a sheet
func (a *App) HandleEditSheet(vr *views.ViewRenderer) error {
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
	return vr.EditSheetPage(dateStr, content)
}

// HandleGetSheetByDate handles GET /sheets/{date} - returns a specific sheet by date
func (a *App) HandleGetSheetByDate(vr *views.ViewRenderer) error {
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
	return vr.SheetListingComponent(current)
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
	return vr.SheetListingComponent(sheet)
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
	return vr.SheetListingComponent(sheet)
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

// HandleGetChangePattern handles GET /change-pattern - shows the pattern editor
func (a *App) HandleGetChangePattern(vr *views.ViewRenderer) error {
	pattern := a.sheetService.GetPattern()
	// Convert pattern to a map for easy lookup
	selectedMap := make(map[int]bool)
	for _, day := range pattern {
		if day >= 1 && day <= 200 {
			selectedMap[day] = true
		}
	}
	return vr.ChangePatternPage(selectedMap)
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

	// Search for matching sheets
	results, err := a.sheetService.Search(query)
	if err != nil {
		return err
	}

	return vr.SheetListingComponents(results)
}
