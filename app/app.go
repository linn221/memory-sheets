package app

import (
	"fmt"
	"net/http"
	"os"
)

type App struct {
	sheetService    *SheetService
	navSheetService *NavSheetService
	AuthMiddleware  func(http.Handler) http.Handler
	config          Cfg
}

type Cfg struct {
	port string
}

func NewApp(dir string, port string, authMiddleware func(http.Handler) http.Handler, pattern RemindPattern) *App {
	// Try to load pattern from JSON file, fallback to provided pattern
	patternFile := "pattern.json"
	loadedPattern, err := LoadPatternFromJSON(patternFile)
	if err != nil {
		// If loading fails, use provided pattern and try to save it
		loadedPattern = pattern
		if saveErr := SavePatternToJSON(patternFile, pattern); saveErr != nil {
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
		dir: "nav",
	}
	// Create nav directory if it doesn't exist
	if err := os.MkdirAll("nav", 0755); err != nil {
		fmt.Printf("Warning: failed to create nav directory: %v\n", err)
	}
	err = navSheetService.ReadDir()
	if err != nil {
		// If nav directory is empty or has issues, that's okay
		fmt.Printf("Warning: failed to read nav directory: %v\n", err)
	}

	return &App{
		config: Cfg{
			port: port,
		},
		sheetService:    sheetSerice,
		navSheetService: navSheetService,
		AuthMiddleware:  authMiddleware,
	}
}
