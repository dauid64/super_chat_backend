package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"

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
	fromUserID := r.URL.Query().Get("fromuser")
	toUserID := r.URL.Query().Get("touser")

	var messagesFromUser []models.Message
	var messsagesToUser []models.Message

	var waitGroup sync.WaitGroup
	waitGroup.Add(2)

	go func() {
		record := database.Instance.Joins("ToUser").Joins("FromUser").Find(&messagesFromUser, "messages.from_user_id = ?", fromUserID)
		if record.Error != nil {
			responses.Erro(w, http.StatusInternalServerError, record.Error)
		}
		waitGroup.Done()
	}()

	go func() {
		record := database.Instance.Joins("FromUser").Joins("ToUser").Find(&messsagesToUser, "messages.to_user_id = ?", toUserID)
		if record.Error != nil {
			responses.Erro(w, http.StatusInternalServerError, record.Error)
		}
		waitGroup.Done()
	}()

	waitGroup.Wait()

	type messagesData struct {
		FromMessages []models.Message `json:"fromMessages,omitempty"`
		ToMessages []models.Message	`json:"toMessages,omitempty"`
	}
	messagesResponse := messagesData{
		messagesFromUser,
		messsagesToUser,
	}

	responses.JSON(w, http.StatusOK, messagesResponse)
}