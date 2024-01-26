package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dauid64/super_chat_backend/src/config"
	"github.com/dauid64/super_chat_backend/src/database"
	"github.com/dauid64/super_chat_backend/src/router"
)

func main() {
	config.LoadEnvironment()

	database.Conect()
	database.Migrate()

	r := router.Generate()
	log.Printf("Escutando na porta %d\n", config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r))
}