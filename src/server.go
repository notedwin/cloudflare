package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/lestrrat-go/jwx/jwt"
)

var tokenAuth *jwtauth.JWTAuth

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("token")
			ctoken := ""
			if err != nil {
				ctoken = cookie.Value
			}

			token, err := tokenAuth.Decode(ctoken)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
			}
			if token == nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
			}
			if err := jwt.Validate(token); err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
			}

			var claims map[string]interface{} = token.PrivateClaims()

			w.Write([]byte(fmt.Sprintf("protected area. hi %v", claims["sub"])))
		})
	})
	r.Get("/auth/{username}", func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")

		_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"sub": username, "exp": 86400})

		//return token as a cookie
		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: time.Now().Add(time.Hour * 24),
		})
		w.Write([]byte(fmt.Sprintf("Welcome %v", username)))
	})

	r.Get("/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	//serve README.txt
	FileServer(r, "/", http.Dir("./"))

	//automatically returns 404 for missing routes
	http.ListenAndServe(":3000", r)
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
