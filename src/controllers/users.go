package controllers

import (
	"net/http"

	"github.com/dauid64/super_chat_backend/src/database"
	"github.com/dauid64/super_chat_backend/src/models"
	"github.com/dauid64/super_chat_backend/src/responses"
)

func SearchUsers(w http.ResponseWriter, r *http.Request) {
	db, err := database.Conect()

	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	lines, err := db.Query("SELECT id, email, password FROM users")

	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		return
	}

	var users []models.User

	for lines.Next() {
		var user models.User

		err = lines.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
		)

		if err != nil {
			responses.Erro(w, http.StatusInternalServerError, err)
			return
		}

		users = append(users, user)
	}

	responses.JSON(w, http.StatusOK, users)
}