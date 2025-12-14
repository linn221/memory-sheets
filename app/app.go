package app

import (
	"fmt"
	"net/http"
)

type App struct {
	sheetService   *SheetService
	AuthMiddleware func(http.Handler) http.Handler
	config         Cfg
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
	return &App{
		config: Cfg{
			port: port,
		},
		sheetService:   sheetSerice,
		AuthMiddleware: authMiddleware,
	}
}
