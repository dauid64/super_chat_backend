package routers

import (
	"net/http"

	"github.com/dauid64/super_chat_backend/src/controllers"
)

var MessagesRouters = []Router{
	{
		URI: "/mensagens",
		Method: http.MethodPost,
		Func: controllers.CreateMessage,
		isAutenticate: true,
	},
	{
		URI: "/mensagens",
		Method: http.MethodGet,
		Func: controllers.GetMessagesChat,
		isAutenticate: true,
	},
}