package app

import (
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

	sheetSerice := &SheetService{
		pattern: pattern,
		dir:     dir,
	}
	err := sheetSerice.ReadDir()
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
