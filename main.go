package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dauid64/super_chat_backend/src/config"
	"github.com/dauid64/super_chat_backend/src/database"
	"github.com/dauid64/super_chat_backend/src/router"
	"github.com/rs/cors"
)

func main() {
	config.LoadEnvironment()

	database.Conect()
	database.Migrate()

	r := router.Generate()

	cors := cors.New(cors.Options{
        AllowedOrigins:   []string{config.FrontEndUrl},
        AllowedHeaders:   []string{"*"},
        AllowedMethods:   []string{"GET", "PATCH", "POST", "PUT", "OPTIONS", "DELETE"},
        AllowCredentials: true,
    })

	handler := cors.Handler(r)

	log.Printf("Escutando na porta %d\n", config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), handler))
}