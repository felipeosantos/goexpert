package clients_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/felipeosantos/curso-go/goexpert/cloud-run/internal/infra/webserver/clients"
	"github.com/stretchr/testify/assert"
)

// MockHTTPClientWeatherAPI implements HTTPClient interface for testing
type MockHTTPClientWeatherAPI struct {
	Do func(req *http.Request) (*http.Response, error)
}

// Satisfy the interface expected by BuscaCurrentWeather
func (m *MockHTTPClientWeatherAPI) DoFunc(req *http.Request) (*http.Response, error) {
	return m.Do(req)
}

func TestBuscaCurrentWeather_Success(t *testing.T) {
	os.Setenv("WEATHER_API_KEY", "dummykey")
	defer os.Unsetenv("WEATHER_API_KEY")

	mockResponse := `{
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
	mockClient := &MockHTTPClientWeatherAPI{
		Do: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(mockResponse)),
			}, nil
		},
	}

	result, weatherErr, err := clients.BuscaCurrentWeather(context.Background(), mockClient, "Guarulhos")
	assert.NoError(t, err)
	assert.Nil(t, weatherErr)
	assert.NotNil(t, result)
	assert.Equal(t, "Guarulhos", result.Location.Name)
	assert.Equal(t, "21.1", result.Current.TempC)
}

func TestBuscaCurrentWeather_HTTPClientError(t *testing.T) {
	os.Setenv("WEATHER_API_KEY", "dummykey")
	defer os.Unsetenv("WEATHER_API_KEY")

	mockClient := &MockHTTPClientWeatherAPI{
		Do: func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("network error")
		},
	}

	result, weatherErr, err := clients.BuscaCurrentWeather(context.Background(), mockClient, "Guarulhos")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, weatherErr)
}

func TestBuscaCurrentWeather_StatusCodeError(t *testing.T) {
	os.Setenv("WEATHER_API_KEY", "dummykey")
	defer os.Unsetenv("WEATHER_API_KEY")

	mockResponse := `{
	    "error": {
	        "code": 1006,
	        "message": "No matching location found."
	    }
	}`
	mockClient := &MockHTTPClientWeatherAPI{
		Do: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(bytes.NewBufferString(mockResponse)),
			}, nil
		},
	}

	result, weatherErr, err := clients.BuscaCurrentWeather(context.Background(), mockClient, "UnknownPlace")
	assert.NoError(t, err)
	assert.NotNil(t, weatherErr)
	assert.Nil(t, result)
	assert.Equal(t, 1006, weatherErr.Error.Code)
	assert.Equal(t, "No matching location found.", weatherErr.Error.Message)
	assert.Equal(t, http.StatusBadRequest, weatherErr.HttpStatusCode)
}
