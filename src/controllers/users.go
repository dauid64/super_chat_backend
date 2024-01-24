package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dauid64/super_chat_backend/src/database"
)

func SearchUsers(w http.ResponseWriter, r *http.Request) {
	_, err := database.Conect()
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)

		err := json.NewEncoder(w).Encode(struct {
			Erro string `json:"erro"`
		}{
			Erro: err.Error(),
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	err = json.NewEncoder(w).Encode("Conectado!")
}