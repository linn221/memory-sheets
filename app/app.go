package app

import (
	"net/http"

	"github.com/linn221/memory-sheets/services"
)

type App struct {
	SheetService   *services.SheetService
	AuthMiddleware func(http.Handler) http.Handler
	port           string
}

func NewApp(dir string, port string, authMiddleware func(http.Handler) http.Handler) *App {
	sheetService := services.NewSheetService(dir)
	return &App{
		SheetService:   sheetService,
		AuthMiddleware: authMiddleware,
		port:           port,
	}
}
