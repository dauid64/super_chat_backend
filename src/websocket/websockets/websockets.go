package websockets

import (
	"net/http"

	"github.com/gorilla/mux"
)

type WebSocket struct {
	URI string
	Func func(http.ResponseWriter, *http.Request)
}

func Configurate(r *mux.Router) *mux.Router {
	go broadcaster()
	
	websockets := ChatWebSocket

	for _, websocket := range websockets {
		r.HandleFunc("/ws" + websocket.URI, websocket.Func)
	}

	return r
}