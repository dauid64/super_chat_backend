package controllers

import (
	"encoding/json"
	"log"
	"net/http"
)

func SearchUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	err := json.NewEncoder(w).Encode("Olá Mundo")
	if err != nil {
		log.Fatal(err)
	}
}