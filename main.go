package main

import (
	"github.com/linn221/memory-sheets/app"
	secretmiddleware "github.com/linn221/memory-sheets/secretMiddleware"
)

func main() {
	secretMd := secretmiddleware.New("http://localhost", "8033", "/secret", "/", secretmiddleware.PersistentSecret("secret.txt"))
	app := app.NewApp("sheets", "8033", secretMd, app.RemindPattern{1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89})
	app.Serve()
}
