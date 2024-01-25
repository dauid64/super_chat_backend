package responses

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dauid64/super_chat_backend/src/config"
)


func JSON(w http.ResponseWriter, statusCode int, dados interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if dados != nil {
		err := json.NewEncoder(w).Encode(dados)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func Erro(w http.ResponseWriter, statusCode int, err error) {
	JSON(w, statusCode, struct{
		Erro string `json:"erro"`
	}{
		Erro: err.Error(),
	})
}

func EnableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", config.FrontEndUrl)
}