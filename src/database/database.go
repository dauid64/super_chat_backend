package database

import (
	"database/sql"

	"github.com/dauid64/super_chat_backend/src/config"
)

// Conectar abre a conexão com o banco de dados e a retorna
func Conect() (*sql.DB, error) {
	db, err := sql.Open("postgres", config.DataBaseSourceName)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}