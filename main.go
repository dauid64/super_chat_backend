package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dauid64/super_chat_backend/src/config"
	"github.com/dauid64/super_chat_backend/src/router"
)

func main() {
	config.Carregar()
	r := router.Generate()

	fmt.Printf("Escutando na porta %d\n", config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r))
}