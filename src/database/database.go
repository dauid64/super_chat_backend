package database

import (
	"log"
	"testing"

	"github.com/dauid64/super_chat_backend/src/config"
	"github.com/dauid64/super_chat_backend/src/models"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Instance *gorm.DB
var dbError error

// Conectar abre a conexão com o banco de dados e a retorna
func Conect() {
	Instance, dbError = gorm.Open(postgres.Open(config.DataBaseSourceName), &gorm.Config{TranslateError: true})
	if dbError != nil {
		log.Fatal(dbError)
	}

	log.Println("Connected to Database!")
}

func Migrate() {
	Instance.AutoMigrate(&models.User{}, &models.Message{})
	log.Println("Database Migrations Completed!")
}

func DbMock(t *testing.T) sqlmock.Sqlmock {
	sqldb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	gormdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqldb,
	}))

	if err != nil {
		t.Fatal(err)
	}

	Instance = gormdb

	return mock
}