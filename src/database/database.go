package database

import (
	"log"

	"github.com/dauid64/super_chat_backend/src/config"
	"github.com/dauid64/super_chat_backend/src/models"
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