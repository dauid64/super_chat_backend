package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/dauid64/super_chat_backend/src/authetication"
	"github.com/dauid64/super_chat_backend/src/database"
	"github.com/dauid64/super_chat_backend/src/models"
	"github.com/dauid64/super_chat_backend/src/responses"
	"github.com/gorilla/mux"
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

func GetMessagesChat(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)

	usuarioID, err := authetication.ExtractUserID(r)
	if err != nil {
		responses.Erro(w, http.StatusUnauthorized, err)
		return
	}

	toUserID, err := strconv.ParseUint(param["touser"], 10, 64)
	if err != nil {
		responses.Erro(w, http.StatusBadRequest, err)
		return
	}

	var messages []models.Message

	record := database.Instance.Joins(
		"ToUser").Joins("FromUser").Where(
		"messages.from_user_id IN ? AND messages.to_user_id IN ?", []uint64{usuarioID, toUserID}, []uint64{usuarioID, toUserID},
	).Order("created_at ASC").Find(&messages)
	if record.Error != nil {
		responses.Erro(w, http.StatusInternalServerError, record.Error)
	}

	responses.JSON(w, http.StatusOK, messages)
}
