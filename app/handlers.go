package app

import (
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/linn221/memory-sheets/models"
	"github.com/linn221/memory-sheets/views"
)

func ActionHandler(httpFunc func(w http.ResponseWriter), handle func(session *models.Session, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := TheSession
		err := handle(&session, r)
		if err != nil {
			views.ErrorBox(err.Error()).Render(r.Context(), w)
			return
		}
		httpFunc(w)
	}
}

// HandleGetSheets handles GET /sheets - returns sheets that need to be reminded today
func (a *App) HandleGetSheets(vr *views.ViewRenderer) error {
	today := Today()
	remindingSheets, err := a.sheetService.LookUpSheets(today)
	if err != nil {
		return err
	}
	return vr.ListSheets(&TheSession, remindingSheets)
}

var RedirectIndex = func(w http.ResponseWriter) {
	w.Header().Set("HX-Location", "/sheets")
	w.WriteHeader(http.StatusOK)
}

// HandleGetAllSheets handles GET /all-sheets - returns all sheets
func (a *App) HandleGetAllSheets(vr *views.ViewRenderer) error {
	return vr.ListSheets(&TheSession, a.sheetService.sheets)
}

// HandleCreateSheet handles POST /sheets - creates a new sheet for today
func (a *App) HandleCreateSheet() http.HandlerFunc {
	return ActionHandler(RedirectIndex, func(session *models.Session, r *http.Request) error {
		today := Today()

		// Read request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return errors.New("failed to read request body")
		}
		defer r.Body.Close()

		content := string(body)

		// Check if sheet already exists
		if a.sheetService.IsSheetExist(today) {
			return errors.New("sheet already exists for today")
		}

		// Create the sheet
		err = a.sheetService.CreateSheet(today, content)
		if err != nil {
			return err
		}

		session.SetMessage("sheet created successfully")

		return nil
	})
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
	var current, prev, next *models.MemorySheet
	for i, rSheet := range remindingSheets {
		if rSheet.Date.Equal(date) {
			current = rSheet
			if i > 0 {
				prev = remindingSheets[i-1]
			}
			if i+1 != len(remindingSheets) {
				next = remindingSheets[i+1]
			}
			break
		}
	}
	if current == nil {
		return errors.New("note not found")
	}
	return vr.ListingSheet(current, prev, next)
}

// HandleUpdateSheet handles PUT /sheets/{date} - updates an existing sheet
func (a *App) HandleUpdateSheet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dateStr := r.PathValue("date")
		redirectToDate := func(w http.ResponseWriter) {
			w.Header().Set("HX-Location", "/sheets/"+dateStr)
			w.WriteHeader(http.StatusOK)
		}
		ActionHandler(redirectToDate, func(session *models.Session, r *http.Request) error {
			if dateStr == "" {
				return errors.New("date cannot be empty")
			}

			date, err := time.Parse(time.DateOnly, dateStr)
			if err != nil {
				return errors.New("invalid date format")
			}

			// Read request body
			body, err := io.ReadAll(r.Body)
			if err != nil {
				return errors.New("failed to read request body")
			}
			defer r.Body.Close()

			content := string(body)

			// Update the sheet
			err = a.sheetService.UpdateSheet(date, content)
			if err != nil {
				return err
			}

			return nil
		})(w, r)
	}
}

// HandleDeleteSheet handles DELETE /sheets/{date} - deletes a sheet
func (a *App) HandleDeleteSheet() http.HandlerFunc {
	return ActionHandler(RedirectIndex, func(session *models.Session, r *http.Request) error {
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
	})
}
