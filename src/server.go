package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt"
)

var hmacSampleSecret []byte = []byte("secret")

//create a new token
func createToken(username string) (string, error) {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(hmacSampleSecret)
	return tokenString, err
}

func validateToken(tokenString string) (bool, map[string]interface{}, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSampleSecret, nil
	})
	if err != nil {
		return false, map[string]interface{}{}, err
	} else {
		//get claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			fmt.Printf("%v %v\n", claims["exp"], claims["iat"])
			return true, claims, nil
		}
		return false, map[string]interface{}{"error": err.Error()}, nil
	}

}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie("token")
		if cookie != nil {
			valid, claims, err := validateToken(cookie.Value)
			if valid {
				w.Write([]byte(fmt.Sprintf("Welcome to the admin page, %v!", claims["sub"])))
			} else {
				w.Write([]byte(fmt.Sprintf("Invalid Token: %v", err)))
			}
		} else {
			w.Write([]byte("Please login\n"))
		}
	})

	r.Get("/auth/{username}", func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")

		tokenString, err := createToken(username)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		//return token as a cookie that is not accessible by javascript
		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Path:    "/",
			Expires: time.Now().Add(time.Hour * 24),
			Secure:  true,
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
