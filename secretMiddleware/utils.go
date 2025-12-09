package secretmiddleware

import (
	"math/rand"
	"net/http"
	"time"
)

func removeCookies(w http.ResponseWriter, key string) {
	http.SetCookie(w, &http.Cookie{
		Name:    key,
		Expires: time.Unix(0, 0), // Set to past
		MaxAge:  -1,              // Also ensures deletion
		Path:    "/",
		Domain:  "",
	})
}

// set secure cookies
func setCookies(w http.ResponseWriter, key string, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:   key,
		Value:  value,
		MaxAge: 0,
		Path:   "/", Domain: "",
		Secure: false, HttpOnly: true,
	})
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomString(length int) string {
	// rand.Seed(time.Now().UnixNano())
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
