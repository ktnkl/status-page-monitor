package middleware

import (
	"log"
	"net/http"
)

func Logging() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.URL.Path)
			f(w, r)
		}
	}
}
