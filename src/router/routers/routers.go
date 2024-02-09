package routers

import (
	"net/http"

	"github.com/dauid64/super_chat_backend/src/middlewares"
	"github.com/dauid64/super_chat_backend/src/websocket/websockets"

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
	routers = append(routers, LoginRouter)
	routers = append(routers, MessagesRouters...)
	
	for _, router := range routers {
		if router.isAutenticate {
			r.HandleFunc(router.URI, middlewares.CorsHandler(middlewares.Logger(middlewares.Authenticate(router.Func))),
			).Methods(router.Method)
		} else {
			r.HandleFunc(router.URI, middlewares.CorsHandler(middlewares.Logger(router.Func))).Methods(router.Method)
		}
	}

	hubs := make( map[string]*websockets.Hub )

	r.HandleFunc("/ws/chat/{fromUserID}/{toUserID}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		toUserID := params["toUserID"]

		hub, ok := hubs[toUserID]
		if !ok {
			hubs[toUserID] = websockets.NewHub()
			hub = hubs[toUserID]
			hub.Run()
		}

		websockets.ServeWs(hub, w, r)
	})

	return r
}
