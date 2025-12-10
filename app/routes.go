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

	mux.HandleFunc("GET /sheets", views.Handler(func(vr *views.ViewRenderer) error {
		today := Today()
		remindingSheets, err := a.sheetService.LookUpSheets(a.storage.sheets, today, a.storage.pattern)
		if err != nil {
			return err
		}
		return vr.ListSheets(remindingSheets)
	}))
	mux.HandleFunc("POST /sheets", func(w http.ResponseWriter, r *http.Request) {

	})
	mux.HandleFunc("GET /sheets/{date}", func(w http.ResponseWriter, r *http.Request) {

	})
	mux.HandleFunc("PUT /sheets/{date}", func(w http.ResponseWriter, r *http.Request) {

	})
	mux.HandleFunc("DELETE /sheets/{date}", func(w http.ResponseWriter, r *http.Request) {

	})
}
