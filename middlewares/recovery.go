package middlewares

import (
	"log"
	"net/http"
	"runtime/debug"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {

				// Log the panic message and stack trace
				log.Printf("[PANIC RECOVERED] %v\n%s", rec, debug.Stack())

				// Optional: customize the error response
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
