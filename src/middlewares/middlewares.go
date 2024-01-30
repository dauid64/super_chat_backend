package middlewares

import (
	"log"
	"net/http"

	"github.com/dauid64/super_chat_backend/src/authetication"
	"github.com/dauid64/super_chat_backend/src/config"
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

func CorsHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if (r.Method == "OPTIONS") {
			log.Print("preflight detected: ", r.Header)
			w.Header().Add("Connection", "keep-alive") 			
			w.Header().Add("Access-Control-Allow-Origin", "http://localhost:3000") 			
			w.Header().Add("Access-Control-Allow-Methods", "POST, OPTIONS, GET, DELETE, PUT") 			
			w.Header().Add("Access-Control-Allow-Headers", "content-type") 			
			w.Header().Add("Access-Control-Max-Age", "86400")
			w.WriteHeader(http.StatusOK)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", config.FrontEndUrl)
			w.Header().Set("Access-Control-Allow-Headers", "authorization,content-type")
			w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, PATCH, GET, DELETE, OPTIONS")
			next(w, r)
		}	
	}
}