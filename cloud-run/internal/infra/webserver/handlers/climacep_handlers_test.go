package handlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/felipeosantos/curso-go/goexpert/cloud-run/internal/entity"
	"github.com/felipeosantos/curso-go/goexpert/cloud-run/internal/infra/webserver/handlers"
	"github.com/stretchr/testify/assert"
)

// MockHTTPClientViaCEP implements HTTPClient for testing
type MockHTTPClientViaCEP struct {
	Do func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClientViaCEP) DoFunc(req *http.Request) (*http.Response, error) {
	return m.Do(req)
}

// MockHTTPClientWeatherAPI implements HTTPClient interface for testing
type MockHTTPClientWeatherAPI struct {
	Do func(req *http.Request) (*http.Response, error)
}

// Satisfy the interface expected by BuscaCurrentWeather
func (m *MockHTTPClientWeatherAPI) DoFunc(req *http.Request) (*http.Response, error) {
	return m.Do(req)
}

func TestBuscaClimaCEP_Success(t *testing.T) {

	mockResponseViaCEP := `{
	  "cep": "01001-000",
	  "logradouro": "Praça da Sé",
	  "complemento": "lado ímpar",
	  "unidade": "",
	  "bairro": "Sé",
	  "localidade": "São Paulo",
	  "uf": "SP",
	  "estado": "São Paulo",
	  "regiao": "Sudeste",
	  "ibge": "3550308",
	  "gia": "1004",
	  "ddd": "11",
	  "siafi": "7107"
	}`

	mockClientViaCEP := &MockHTTPClientViaCEP{
		Do: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(mockResponseViaCEP)),
			}, nil
		},
	}

	os.Setenv("WEATHER_API_KEY", "dummykey")
	defer os.Unsetenv("WEATHER_API_KEY")

	mockResponseWeather := `{
	    "location": {
	        "name": "Guarulhos",
	        "region": "Sao Paulo",
	        "country": "Brazil",
	        "lat": -23.4667,
	        "lon": -46.5333,
	        "tz_id": "America/Sao_Paulo",
	        "localtime_epoch": 1753110612,
	        "localtime": "2025-07-21 12:10"
	    },
	    "current": {
	        "last_updated_epoch": 1753110000,
	        "last_updated": "2025-07-21 12:00",
	        "temp_c": 21.1,
	        "temp_f": 70.0,
	        "is_day": 1,
	        "condition": {
	            "text": "Parcialmente nublado",
	            "icon": "//cdn.weatherapi.com/weather/64x64/day/116.png",
	            "code": 1003
	        },
	        "wind_mph": 4.5,
	        "wind_kph": 7.2,
	        "wind_degree": 222,
	        "wind_dir": "SW",
	        "pressure_mb": 1020.0,
	        "pressure_in": 30.12,
	        "precip_mm": 0.0,
	        "precip_in": 0.0,
	        "humidity": 64,
	        "cloud": 50,
	        "feelslike_c": 21.1,
	        "feelslike_f": 70.0,
	        "windchill_c": 24.3,
	        "windchill_f": 75.7,
	        "heatindex_c": 24.8,
	        "heatindex_f": 76.6,
	        "dewpoint_c": 8.8,
	        "dewpoint_f": 47.8,
	        "vis_km": 10.0,
	        "vis_miles": 6.0,
	        "uv": 5.5,
	        "gust_mph": 5.1,
	        "gust_kph": 8.3
	    }
	}`
	mockClientWeather := &MockHTTPClientWeatherAPI{
		Do: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(mockResponseWeather)),
			}, nil
		},
	}

	handler := &handlers.ClimaCEPHandler{
		ViaCEPHTTPClient:     mockClientViaCEP,
		WeatherAPIHTTPClient: mockClientWeather,
	}

	req := httptest.NewRequest("GET", "/?cep=01001000", nil)
	w := httptest.NewRecorder()

	handler.BuscaClimaCEP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp entity.ClimaResult
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "Guarulhos", resp.City)
	assert.Equal(t, 21.1, resp.TempCelsius)
	assert.Equal(t, 69.98, resp.TempFahrenheit)
	assert.Equal(t, 294.25, resp.TempKelvin)
}

func TestBuscaClimaCEP_InvalidCEP(t *testing.T) {

	mockClientViaCEP := &MockHTTPClientViaCEP{
		Do: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("")),
			}, nil
		},
	}

	os.Setenv("WEATHER_API_KEY", "dummykey")
	defer os.Unsetenv("WEATHER_API_KEY")

	mockClientWeather := &MockHTTPClientWeatherAPI{
		Do: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("")),
			}, nil
		},
	}

	handler := &handlers.ClimaCEPHandler{
		ViaCEPHTTPClient:     mockClientViaCEP,
		WeatherAPIHTTPClient: mockClientWeather,
	}

	req := httptest.NewRequest("GET", "/?cep=0700000A", nil)
	w := httptest.NewRecorder()

	handler.BuscaClimaCEP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	var resp entity.ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "invalid zipcode", resp.Mensagem)
}

func TestBuscaClimaCEP_NotFoundCEP(t *testing.T) {

	mockResponseViaCEP := `{
	  "erro": "true"
	}`

	mockClientViaCEP := &MockHTTPClientViaCEP{
		Do: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(mockResponseViaCEP)),
			}, nil
		},
	}

	os.Setenv("WEATHER_API_KEY", "dummykey")
	defer os.Unsetenv("WEATHER_API_KEY")

	mockClientWeather := &MockHTTPClientWeatherAPI{
		Do: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("")),
			}, nil
		},
	}

	handler := &handlers.ClimaCEPHandler{
		ViaCEPHTTPClient:     mockClientViaCEP,
		WeatherAPIHTTPClient: mockClientWeather,
	}

	req := httptest.NewRequest("GET", "/?cep=07000000", nil)
	w := httptest.NewRecorder()

	handler.BuscaClimaCEP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var resp entity.ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "can not find zipcode", resp.Mensagem)
}
