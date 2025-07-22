package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/felipeosantos/curso-go/goexpert/cloud-run/internal/entity"
	"github.com/felipeosantos/curso-go/goexpert/cloud-run/internal/infra/webserver/clients"
)

type ClimaCEPHandler struct {
	ViaCEPHTTPClient     clients.HTTPClient
	WeatherAPIHTTPClient clients.HTTPClient
}

func NewClimaCEPHandler(viaCEPClient, weatherAPIClient clients.HTTPClient) *ClimaCEPHandler {
	return &ClimaCEPHandler{
		ViaCEPHTTPClient:     viaCEPClient,
		WeatherAPIHTTPClient: weatherAPIClient,
	}
}

func (cch *ClimaCEPHandler) BuscaClimaCEP(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	cepParam := r.URL.Query().Get("cep")

	isValid := isValidCEP(cepParam)

	if !isValid {
		w.WriteHeader(http.StatusUnprocessableEntity)
		errorResponse := entity.ErrorResponse{
			Mensagem: "invalid zipcode",
		}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	cep, err := clients.BuscaCEP(cch.ViaCEPHTTPClient, cepParam)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := entity.ErrorResponse{
			Mensagem: "internal server error",
		}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Verificar se a resposta contém um erro da API (ViaCEP retorna "erro": true para CEPs inválidos)
	if cep.Erro == "true" {
		w.WriteHeader(http.StatusNotFound)
		errorResponse := entity.ErrorResponse{
			Mensagem: "can not find zipcode",
		}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Buscar o clima atual usando a localidade do CEP
	currentWeather, weatherError, err := clients.BuscaCurrentWeather(cch.WeatherAPIHTTPClient, cep.Localidade)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := entity.ErrorResponse{
			Mensagem: "internal server error",
		}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	if weatherError != nil {
		w.WriteHeader(weatherError.HttpStatusCode)
		errorResponse := entity.ErrorResponse{
			Mensagem: weatherError.Error.Message,
		}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	// Unificando os dados do CEP e do clima atual
	climaResult := entity.ClimaResult{
		TempCelsius:    currentWeather.Current.TempC,
		TempFahrenheit: celsiusToFahrenheit(currentWeather.Current.TempC),
		TempKelvin:     celsiusToKelvin(currentWeather.Current.TempC),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(climaResult)

}

func isValidCEP(cep string) bool {
	// Validando formato do CEP (formato mais completo)
	isValid := true
	if len(cep) != 8 {
		isValid = false
	}

	// Validando se foi informado numeros
	if _, err := strconv.Atoi(cep); err != nil {
		isValid = false
	}

	return isValid
}

func celsiusToFahrenheit(celsius float64) float64 {
	return (celsius * 9 / 5) + 32
}

func celsiusToKelvin(celsius float64) float64 {
	return celsius + 273.15
}
