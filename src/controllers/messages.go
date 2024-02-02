package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/dauid64/super_chat_backend/src/authetication"
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

func GetMessagesChat(w http.ResponseWriter, r *http.Request) {
	usuarioID, err := authetication.ExtractUserID(r)
	if err != nil {
		responses.Erro(w, http.StatusUnauthorized, err)
		return
	}

	fromUserID, err := strconv.ParseUint(r.URL.Query().Get("fromuser"), 10, 64)
	if err != nil {
		responses.Erro(w, http.StatusBadRequest, err)
		return
	}
	toUserID, err := strconv.ParseUint(r.URL.Query().Get("touser"), 10, 64)
	if err != nil {
		responses.Erro(w, http.StatusBadRequest, err)
		return
	}

	if usuarioID != fromUserID {
		responses.Erro(w, http.StatusForbidden, errors.New("Você não tem permissão para acessar essa conversa"))
		return
	}

	var messages []models.Message

	record := database.Instance.Joins(
		"ToUser").Joins("FromUser").Where(
		"messages.from_user_id IN ? AND messages.to_user_id IN ?", []uint64{fromUserID, toUserID}, []uint64{fromUserID, toUserID},
	).Order("created_at ASC").Find(&messages)
	if record.Error != nil {
		responses.Erro(w, http.StatusInternalServerError, record.Error)
	}

	responses.JSON(w, http.StatusOK, messages)
}
