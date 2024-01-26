package middlewares

import (
	"log"
	"net/http"

	"github.com/dauid64/super_chat_backend/src/authetication"
	"github.com/dauid64/super_chat_backend/src/responses"
)

func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.Host)
		next(w, r)
	}
}

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := authetication.ValidateToken(r); err != nil {
			responses.Erro(w, http.StatusUnauthorized, err)
			return
		}
		next(w, r)
	}
}