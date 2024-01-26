package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

var (
	// Porta onde a API vai estar rodando
	Port = 0

	// String de conexão com o Postgres
	DataBaseSourceName = ""

	// Url autorizada a realizar consultas CORS
	FrontEndUrl = ""

	// SecretKey é a chave que vai ser usada para assinar o token
	SecretKey []byte
)

func LoadEnvironment() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Port, err = strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		Port = 8000
	}

	DataBaseSourceName = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", 
		os.Getenv("POSTGRES_HOST"), 
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"), 
		os.Getenv("POSTGRES_DB"),
	)

	FrontEndUrl = os.Getenv("FRONTEND_URL")

	SecretKey = []byte(os.Getenv("SECRET_KEY"))
}