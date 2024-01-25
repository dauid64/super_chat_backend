package controllers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/dauid64/super_chat_backend/src/database"
	"github.com/dauid64/super_chat_backend/src/models"

	"github.com/dauid64/super_chat_backend/src/repositories"
	"github.com/dauid64/super_chat_backend/src/responses"
)

func SearchUsers(w http.ResponseWriter, r *http.Request) {
	db, err := database.Conect()

	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repositorie := repositories.NewRepositorieOfUsers(db)
	users, err := repositorie.All()
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, users)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
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

	db, err := database.Conect()

	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repositorie := repositories.NewRepositorieOfUsers(db)
	user.ID, err = repositorie.Create(user)
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, user)
}