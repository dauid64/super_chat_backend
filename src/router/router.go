package router

import (
	"github.com/dauid64/super_chat_backend/src/router/routers"
	"github.com/gorilla/mux"
)

func Generate() *mux.Router {
	r := mux.NewRouter()
	return routers.Configurate(r)
}