package controllers

import (
	"net/http"

	"github.com/dauid64/super_chat_backend/src/responses"
)

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, nil)
}