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
	"github.com/dauid64/super_chat_backend/src/security"
)

func Login(w http.ResponseWriter, r *http.Request) {
	responses.EnableCors(&w)

	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		responses.Erro(w, http.StatusUnprocessableEntity, err)
		return
	}

	var user models.User
	if err = json.Unmarshal(bodyRequest, &user); err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	var userSavedDataBase models.User
	record := database.Instance.Where("email = ?", user.Email).First(&userSavedDataBase)
	if record.Error != nil {
		responses.JSON(w, http.StatusInternalServerError, errors.New("Credenciais Inválidas"))
		return
	}

	err = security.CheckPassword(user.Password, userSavedDataBase.Password)
	if err != nil {
		responses.Erro(w, http.StatusUnauthorized, err)
		return
	}

	token, erro := authetication.CreateToken(uint64(userSavedDataBase.ID))
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusOK, models.AuthenticationData{ID: userID, Token: token})
}