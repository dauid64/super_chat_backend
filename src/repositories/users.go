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

func (respositorie Users) Create(user models.User) (uint64, error) {
	statement, err := respositorie.db.Prepare(
		"INSERT INTO users (email, password) VALUES($1, $2) RETURNING id",
	)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	var lastInsertId uint64

	err = statement.QueryRow(user.Email, user.Password).Scan(&lastInsertId)
	if err != nil {
		return 0, err
	}

	return uint64(lastInsertId), nil
}