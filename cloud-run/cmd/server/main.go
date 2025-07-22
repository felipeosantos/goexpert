package main

import (
	"net/http"

	"github.com/felipeosantos/curso-go/goexpert/cloud-run/internal/infra/webserver/clients"
	"github.com/felipeosantos/curso-go/goexpert/cloud-run/internal/infra/webserver/handlers"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load(".env")

	httpClient := clients.NewCustomClient(&http.Client{})

	climaCEPHandler := handlers.NewClimaCEPHandler(httpClient, httpClient)

	http.HandleFunc("/", climaCEPHandler.BuscaClimaCEP)
	http.ListenAndServe(":8080", nil)
}
