package responses

import (
	"encoding/json"
	"log"
	"net/http"
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