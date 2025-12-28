package app

import (
	"net/http"

	"github.com/linn221/memory-sheets/views"
)

func (a *App) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /sheets", views.Handler(a.ShowTodaySheets))
	mux.HandleFunc("GET /sheets/{date}/edit", views.Handler(a.ShowEditSheet))
	mux.HandleFunc("GET /all-sheets", views.Handler(a.ShowAllSheets))
	mux.HandleFunc("GET /sheets/{date}", views.Handler(a.ShowSheet))
	mux.HandleFunc("POST /sheets", views.Handler(a.HandleCreateSheet))
	mux.HandleFunc("PUT /sheets/{date}", views.Handler(a.HandleUpdateSheet))
	mux.HandleFunc("DELETE /sheets/{date}", views.Handler(a.HandleDeleteSheet))
	mux.HandleFunc("GET /search", views.Handler(a.HandleSearch))
	mux.HandleFunc("GET /nav-sheets/{title}", views.Handler(a.ShowNavSheet))
	mux.HandleFunc("GET /nav-sheets/new", views.Handler(a.ShowCreateNavSheet))
	mux.HandleFunc("POST /nav-sheets", views.Handler(a.HandleCreateNavSheet))
	mux.HandleFunc("GET /nav-sheets/{title}/edit", views.Handler(a.ShowEditNavSheet))
	mux.HandleFunc("PUT /nav-sheets/{title}", views.Handler(a.HandleUpdateNavSheet))
	mux.HandleFunc("DELETE /nav-sheets/{title}", views.Handler(a.HandleDeleteNavSheet))
	// mux.HandleFunc("GET /change-pattern", views.Handler(a.ShowChangePattern))
	// mux.HandleFunc("POST /change-pattern", views.Handler(a.HandlePostChangePattern))
}
