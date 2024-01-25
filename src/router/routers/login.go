package routers

import (
	"net/http"

	"github.com/dauid64/super_chat_backend/src/controllers"
)

var LoginRouter = Router{
	URI:    "/login",
	Method: http.MethodGet,
	Func: controllers.Login,
	isAutenticate: false,
}