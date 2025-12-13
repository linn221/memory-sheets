package app

import (
	"net/http"
	"time"

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
	mux.HandleFunc("GET /sheets/new", views.Handler(func(vr *views.ViewRenderer) error {
		return vr.NewSheetPage()
	}))
	mux.HandleFunc("GET /sheets/{date}/edit", views.Handler(func(vr *views.ViewRenderer) error {
		r := vr.Request()
		dateStr := r.PathValue("date")
		date, err := time.Parse(time.DateOnly, dateStr)
		if err != nil {
			return err
		}
		index, err := a.sheetService.GetSheetByDate(date)
		if err != nil {
			return err
		}
		content := a.sheetService.sheets[index].Text
		return vr.EditSheetPage(dateStr, content)
	}))

	mux.HandleFunc("GET /all-sheets", views.Handler(a.HandleGetAllSheets))
	mux.HandleFunc("POST /sheets", a.HandleCreateSheet())
	mux.HandleFunc("GET /sheets/{date}", views.Handler(a.HandleGetSheetByDate))
	mux.HandleFunc("PUT /sheets/{date}", a.HandleUpdateSheet())
	mux.HandleFunc("DELETE /sheets/{date}", a.HandleDeleteSheet())
	fileHandler := http.StripPrefix("/static", http.FileServer(http.Dir("static")))
	mux.Handle("/static/", fileHandler)
}
