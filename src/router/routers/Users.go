package routers

import (
	"net/http"

	"github.com/dauid64/super_chat_backend/src/controllers"
)

var UsersRouters = []Router{
	{
		URI:           "/usuarios",
		Method:        http.MethodGet,
		Func:          controllers.SearchUsers,
		isAutenticate: true,
	},
	{
		URI:           "/usuarios",
		Method:        http.MethodPost,
		Func:          controllers.CreateUser,
		isAutenticate: false,
	},
	{
		URI: "recuperar/usuario",
		Method: http.MethodGet,
		Func: controllers.RecoverUser,
		isAutenticate: true,
	},
}
