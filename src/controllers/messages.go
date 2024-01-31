package controllers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/dauid64/super_chat_backend/src/database"
	"github.com/dauid64/super_chat_backend/src/models"
	"github.com/dauid64/super_chat_backend/src/responses"
)

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Erro(w, http.StatusUnprocessableEntity, err)
		return
	}

	var message models.Message
	err = json.Unmarshal(bodyRequest, &message)
	if err != nil {
		responses.Erro(w, http.StatusBadRequest, err)
		return
	}

	record := database.Instance.Create(&message)
	if record.Error != nil {
		responses.Erro(w, http.StatusInternalServerError, record.Error)
		return
	}

	responses.JSON(w, http.StatusOK, message)
}