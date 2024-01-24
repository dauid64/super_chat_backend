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
		isAutenticate: false,
	},
}
