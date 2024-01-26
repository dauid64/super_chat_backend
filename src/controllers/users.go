package controllers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/dauid64/super_chat_backend/src/database"
	"github.com/dauid64/super_chat_backend/src/models"

	"github.com/dauid64/super_chat_backend/src/responses"
)

func SearchUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User

	record := database.Instance.Find(&users)
	if record.Error != nil {
		responses.Erro(w, http.StatusInternalServerError, record.Error)
		return
	}
	if record.Error != nil {
		responses.Erro(w, http.StatusInternalServerError, record.Error)
		return
	}

	responses.JSON(w, http.StatusOK, users)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	responses.EnableCors(&w)

	BodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Erro(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user models.User
	err = json.Unmarshal(BodyRequest, &user)
	if err != nil {
		responses.Erro(w, http.StatusBadRequest, err)
		return
	}

	err = user.Prepare("cadastro")
	if err != nil {
		responses.Erro(w, http.StatusBadRequest, err)
		return
	}

	record := database.Instance.Create(&user)
	if record.Error != nil {
		responses.Erro(w, http.StatusInternalServerError, record.Error)
		return
	}

	responses.JSON(w, http.StatusCreated, user)
}