package main

import (
	"fmt"
	"net/http"

	"github.com/linn221/memory-sheets/app"
	"github.com/linn221/memory-sheets/middlewares"
	secretmiddleware "github.com/linn221/memory-sheets/secretMiddleware"
)

func main() {
	secretMd := secretmiddleware.New("http://localhost", "8033", "/secret", "/sheets", secretmiddleware.PersistentSecret("secret.txt"), func(magicLink string) {
		// could decide to do either email or print to console
		fmt.Println(magicLink)
	})
	app := app.NewApp("sheets", "pattern.json")
	mux := http.NewServeMux()
	app.SetupRoutes(mux)

	srv := http.Server{
		Addr:    ":8033",
		Handler: secretMd(middlewares.LoggingMiddleware(middlewares.Recovery(mux))),
	}
	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
