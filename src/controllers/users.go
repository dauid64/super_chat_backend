package controllers

import (
	"fmt"
	"net/http"

	"github.com/dauid64/super_chat_backend/src/database"

	"github.com/dauid64/super_chat_backend/src/repositories"
	"github.com/dauid64/super_chat_backend/src/responses"
)

func SearchUsers(w http.ResponseWriter, r *http.Request) {
	db, err := database.Conect()

	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		fmt.Printf("Url: /usuarios, Method: GET, Response: %d\n", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	repositorie := repositories.NewRepositorieOfUsers(db)
	users, err := repositorie.All()
	if err != nil {
		responses.Erro(w, http.StatusInternalServerError, err)
		fmt.Printf("Url: /usuarios, Method: GET, Response: %d\n", http.StatusInternalServerError)
		return
	}
	fmt.Printf("Url: /usuarios, Method: GET, Response: %d\n", http.StatusOK)
	responses.JSON(w, http.StatusOK, users)
}