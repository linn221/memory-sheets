package app

import (
	"fmt"
	"net/http"
	"os"
	"path"
)

type App struct {
	sheetService    *SheetService
	navSheetService *NavSheetService
}

func Handler(mux *http.ServeMux, dir string, patternFile string) http.Handler {
	app := NewApp(dir, patternFile)
	app.SetupRoutes(mux)
	return mux
}

func NewApp(dir string, patternFile string) *App {
	// Try to load pattern from JSON file, fallback to provided pattern
	loadedPattern, err := LoadPatternFromJSON(patternFile)
	if err != nil {
		// If loading fails, use default pattern and try to save it
		loadedPattern = RemindPattern{1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89}
		if saveErr := SavePatternToJSON(patternFile, loadedPattern); saveErr != nil {
			// Log error but don't fail startup
			fmt.Printf("Warning: failed to save pattern to %s: %v\n", patternFile, saveErr)
		}
	}

	sheetSerice := &SheetService{
		pattern: loadedPattern,
		dir:     dir,
	}
	err = sheetSerice.ReadDir()
	if err != nil {
		panic(err)
	}

	navSheetService := &NavSheetService{
		dir: path.Join(dir, "nav"),
	}
	// Create nav directory if it doesn't exist
	if err := os.MkdirAll(navSheetService.dir, 0755); err != nil {
		fmt.Printf("Warning: failed to create nav directory: %v\n", err)
	}
	err = navSheetService.ReadDir()
	if err != nil {
		// If nav directory is empty or has issues, that's okay
		fmt.Printf("Warning: failed to read nav directory: %v\n", err)
	}

	return &App{
		sheetService:    sheetSerice,
		navSheetService: navSheetService,
	}
}
