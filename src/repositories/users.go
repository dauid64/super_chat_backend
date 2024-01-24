package repositories

import (
	"database/sql"

	"github.com/dauid64/super_chat_backend/src/models"
)

type Users struct {
	db *sql.DB
}

func NewRepositorieOfUsers(db *sql.DB) *Users {
	return &Users{db}
}

func (repositorie Users) All() ([]models.User, error){
	lines, err := repositorie.db.Query("SELECT id, email, password FROM users")
	if err != nil {
		return nil, err
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
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}