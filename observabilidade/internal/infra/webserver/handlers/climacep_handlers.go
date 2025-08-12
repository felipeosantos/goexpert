package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/felipeosantos/curso-go/goexpert/observabilidade/internal/entity"
	"github.com/felipeosantos/curso-go/goexpert/observabilidade/internal/infra/webserver/clients"
)

type ClimaCEPHandler struct {
	HTTPClient clients.HTTPClient
}

func NewClimaCEPHandler(httpClient clients.HTTPClient) *ClimaCEPHandler {
	return &ClimaCEPHandler{
		HTTPClient: httpClient,
	}
}

func (cch *ClimaCEPHandler) BuscaClimaCEP(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var cepRequest entity.ClimaCEPRequest
	if err := json.NewDecoder(r.Body).Decode(&cepRequest); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		errorResponse := entity.ClimaCEPError{
			Mensagem: "invalid zipcode",
		}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Buscar o clima atual usando a localidade do CEP
	climacepBResponse, climacepBError, httpStatusCode, err := clients.BuscaClimaCEPB(r.Context(), cch.HTTPClient, cepRequest.Cep)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := entity.ClimaCEPError{
			Mensagem: "internal server error",
		}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	if climacepBError != nil {
		w.WriteHeader(httpStatusCode)
		json.NewEncoder(w).Encode(climacepBError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)

	json.NewEncoder(w).Encode(climacepBResponse)

}
