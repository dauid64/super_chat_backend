package websocket

import (
	"github.com/dauid64/super_chat_backend/src/websocket/websockets"
	"github.com/gorilla/mux"
)

func Generate(r *mux.Router) *mux.Router {
	return websockets.Configurate(r)
}