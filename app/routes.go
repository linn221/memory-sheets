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
		remindingSheets, err := a.sheetService.LookUpSheets(today)
		if err != nil {
			return err
		}
		return vr.ListSheets(remindingSheets)
	}))

	mux.HandleFunc("GET /all-sheets", views.Handler(func(vr *views.ViewRenderer) error {
		return vr.ListSheets(a.sheetService.sheets)
	}))

	mux.HandleFunc("POST /sheets", func(w http.ResponseWriter, r *http.Request) {
		today := Today()
		if a.sheetService.IsSheetExist(today) {

		} else {

		}
	})
	mux.HandleFunc("GET /sheets/{date}", func(w http.ResponseWriter, r *http.Request) {

	})
	mux.HandleFunc("PUT /sheets/{date}", func(w http.ResponseWriter, r *http.Request) {

	})
	mux.HandleFunc("DELETE /sheets/{date}", func(w http.ResponseWriter, r *http.Request) {

	})
}
