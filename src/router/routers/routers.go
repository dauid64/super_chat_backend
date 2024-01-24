package routers

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct {
	URI           string
	Method        string
	Func          func(http.ResponseWriter, *http.Request)
	isAutenticate bool
}

func Configurate(r *mux.Router) *mux.Router {
	routers := UsersRouters
	
	for _, router := range routers {
		r.HandleFunc(router.URI, router.Func).Methods(router.Method)
	}

	return r
}
