package app

import (
	"net/http"

	"github.com/linn221/memory-sheets/views"
)

func (a *App) Serve() {
	mux := http.NewServeMux()
	a.SetupRoutes(mux)

	srv := http.Server{
		Addr:    ":" + a.config.port,
		Handler: a.AuthMiddleware(loggingMiddleware(recovery(mux))),
	}
	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func (a *App) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /sheets", views.Handler(a.HandleGetSheets))
	mux.HandleFunc("GET /sheets/{date}/edit", views.Handler(a.HandleEditSheet))
	mux.HandleFunc("GET /all-sheets", views.Handler(a.HandleGetAllSheets))
	mux.HandleFunc("GET /sheets/{date}", views.Handler(a.HandleGetSheetByDate))
	mux.HandleFunc("POST /sheets", views.Handler(a.HandleCreateSheet))
	mux.HandleFunc("PUT /sheets/{date}", views.Handler(a.HandleUpdateSheet))
	mux.HandleFunc("DELETE /sheets/{date}", views.Handler(a.HandleDeleteSheet))
	mux.HandleFunc("GET /search", views.Handler(a.HandleSearch))
	// mux.HandleFunc("GET /change-pattern", views.Handler(a.HandleGetChangePattern))
	// mux.HandleFunc("POST /change-pattern", views.Handler(a.HandlePostChangePattern))
	fileHandler := http.StripPrefix("/static", http.FileServer(http.Dir("static")))
	mux.Handle("/static/", fileHandler)
}
