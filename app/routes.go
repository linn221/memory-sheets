package app

import "net/http"

func (a *App) Serve() {
	mux := http.NewServeMux()

	srv := http.Server{
		Addr:    ":" + a.port,
		Handler: a.AuthMiddleware(loggingMiddleware(recovery(mux))),
	}
	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
