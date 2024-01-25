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
	"github.com/dauid64/super_chat_backend/src/repositories"
	"github.com/dauid64/super_chat_backend/src/responses"
	"github.com/dauid64/super_chat_backend/src/security"
)

func Login(w http.ResponseWriter, r *http.Request) {
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

	db, err := database.Conect()
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repositorie := repositories.NewRepositorieOfUsers(db)
	userSalvedDataBase, err := repositorie.SearchEmail(user.Email)
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, errors.New("Esse e-mail não foi cadastrado"))
		return
	}

	if err = security.CheckPassword(user.Password, userSalvedDataBase.Password); err != nil {
		responses.Erro(w, http.StatusInternalServerError, errors.New("Senha inválida"))
		return
	}

	token, err := authetication.CreateToken(userSalvedDataBase.ID)
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	userID := strconv.FormatUint(userSalvedDataBase.ID, 10)

	responses.JSON(w, http.StatusOK, models.AuthenticationData{ID: userID, Token: token})
}