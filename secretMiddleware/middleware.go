package secretmiddleware

import (
	"fmt"
	"net/http"
	"os"
)

// Host: host portion(localhost)
// SecretPath: uri without the slash (start-session)
// RedirectUrl: absolute url to redirect
type SecretConfig struct {
	SecretFunc  func() string
	SecretPath  string
	RedirectUrl string
	Host        string
	PromptFunc  func(string)
	// expiration time.Time // for later

}

func (cfg *SecretConfig) Middleware() func(h http.Handler) http.Handler {
	theSecret := cfg.SecretFunc()
	// fmt.Printf("Magic auth link: %s/%s?secret=%s\n", cfg.Host, cfg.SecretPath, theSecret)
	cfg.PromptFunc(fmt.Sprintf("%s/%s?secret=%s", cfg.Host, cfg.SecretPath, theSecret))

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			currentUrl := r.URL.Path
			if currentUrl == "/"+cfg.SecretPath {
				pretendSecret := r.URL.Query().Get("secret")
				if pretendSecret == theSecret {
					// authentication success
					setCookies(w, "secret", theSecret)
					http.Redirect(w, r, cfg.RedirectUrl, http.StatusTemporaryRedirect)
					return
				}
				http.Error(w, "please visit the magic link for auth", http.StatusUnauthorized)
				return
			}
			cookies, err := r.Cookie("secret")
			if err != nil {
				if err == http.ErrNoCookie {
					http.Error(w, "please visit the magic link for auth", http.StatusUnauthorized)
					return
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			pretendSecret := cookies.Value
			if pretendSecret == theSecret {
				h.ServeHTTP(w, r)
				return
			}

			removeCookies(w, "secret")
			http.Error(w, "please visit the magic link for auth", http.StatusUnauthorized)
		})
	}
}

func New(host string, port string, secretPath string, redirectPath string, secretFunc func() string, promptFunc func(string)) func(http.Handler) http.Handler {
	if redirectPath[0] == '/' {
		redirectPath = redirectPath[1:]
	}
	if secretPath[0] == '/' {
		secretPath = secretPath[1:]
	}
	secretConfig := SecretConfig{
		// Host:        "http://localhost:" + port,
		Host:        host + ":" + port,
		SecretPath:  secretPath,
		RedirectUrl: host + ":" + port + "/" + redirectPath,
		SecretFunc:  secretFunc,
		PromptFunc:  promptFunc,
		// SecretFunc: func() string {
		// 	// return utils.GenerateRandomString(20)
		// 	secretFilename := "secret.txt"
		// 	bs, err := os.ReadFile(secretFilename)
		// 	if err != nil {
		// 		secret := utils.GenerateRandomString(20)
		// 		err := os.WriteFile(secretFilename, []byte(secret), 0666)
		// 		if err != nil {
		// 			panic(err)
		// 		}
		// 		return secret
		// 	}
		// 	return string(bs)
		// },
	}
	return secretConfig.Middleware()
}

var PersistentSecret = func(secretFilename string) func() string {
	return func() string {
		bs, err := os.ReadFile(secretFilename)
		if err != nil {
			secret := generateRandomString(20)
			err := os.WriteFile(secretFilename, []byte(secret), 0666)
			if err != nil {
				panic(err)
			}
			return secret
		}
		return string(bs)
	}
}

var TemporySecret = func() string {
	return generateRandomString(20)
}
