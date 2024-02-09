package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/dauid64/super_chat_backend/src/authetication"
	"github.com/dauid64/super_chat_backend/src/database"
	"github.com/dauid64/super_chat_backend/src/models"
	"github.com/gorilla/mux"

	"github.com/dauid64/super_chat_backend/src/responses"
)

func SearchUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User

	record := database.Instance.Find(&users)
	if record.Error != nil {
		responses.Erro(w, http.StatusInternalServerError, record.Error)
		return
	}

	responses.JSON(w, http.StatusOK, users)
}

func SearchIDUsers(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)

	userID, err := strconv.ParseUint(param["usuarioID"], 10, 64)
	if err != nil {
		responses.Erro(w, http.StatusBadRequest, err)
		return
	}

	var user models.User

	record := database.Instance.Find(&user, "id = ?", userID)
	if record.Error != nil {
		responses.Erro(w, http.StatusInternalServerError, record.Error)
		return
	}

	responses.JSON(w, http.StatusOK, user)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Erro(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user models.User
	err = json.Unmarshal(bodyRequest, &user)
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

func RecoverUser(w http.ResponseWriter, r *http.Request) {
	userIDInToken, erro := authetication.ExtractUserID(r)
	if erro != nil {
		responses.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	var user models.User
	record := database.Instance.Select("id", "created_at", "email").First(&user, "id = ?", userIDInToken)
	if record.Error != nil {
		responses.Erro(w, http.StatusInternalServerError, record.Error)
		return
	}

	responses.JSON(w, http.StatusOK, user)
}