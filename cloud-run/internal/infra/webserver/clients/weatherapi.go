package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/felipeosantos/curso-go/goexpert/cloud-run/internal/entity"
)

func BuscaCurrentWeather(httpClient HTTPClient, localidade string) (*entity.CurrentWeatherResponse, *entity.WeatherErrorResponse, error) {

	apiKey := os.Getenv("WEATHER_API_KEY") // Certifique-se de definir a variável de ambiente WEATHER_API_KEY com sua chave da API

	baseURL := "http://api.weatherapi.com/v1/current.json"
	params := url.Values{}
	params.Set("key", apiKey)
	params.Set("lang", "pt")
	params.Set("q", localidade)

	url := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Erro ao criar requisição para WeatherAPI: %v\n", err)
		return nil, nil, err
	}

	resp, err := httpClient.DoFunc(req)
	if err != nil {
		fmt.Printf("Erro ao fazer requisição para WeatherAPI: %v\n", err)
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Erro ao ler resposta da WeatherAPI: %v\n", err)
		return nil, nil, err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		var weatherError entity.WeatherErrorResponse
		err = json.Unmarshal(body, &weatherError)
		if err != nil {
			fmt.Printf("Erro ao fazer parse da resposta da WeatherAPI: %v\n", err)
			return nil, nil, err
		}

		weatherError.HttpStatusCode = resp.StatusCode

		return nil, &weatherError, nil
	} else {
		var currentWeather entity.CurrentWeatherResponse

		err = json.Unmarshal(body, &currentWeather)
		if err != nil {
			fmt.Printf("Erro ao fazer parse da resposta da WeatherAPI: %v\n", err)
			return nil, nil, err
		}

		return &currentWeather, nil, nil

	}

}
