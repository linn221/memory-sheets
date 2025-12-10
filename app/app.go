package app

import (
	"net/http"
)

type App struct {
	sheetService   *SheetService
	AuthMiddleware func(http.Handler) http.Handler
	config         Cfg
	storage        Storage
}

type Cfg struct {
	port string
	dir  string
}

type Storage struct {
	pattern RemindPattern
	sheets  []*MemorySheet
}

func NewApp(dir string, port string, authMiddleware func(http.Handler) http.Handler, pattern RemindPattern) *App {

	sheetSerice := &SheetService{}
	sheets, err := sheetSerice.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	return &App{
		config: Cfg{
			dir:  dir,
			port: port,
		},
		AuthMiddleware: authMiddleware,
		storage: Storage{
			pattern: pattern,
			sheets:  sheets,
		},
	}
}
