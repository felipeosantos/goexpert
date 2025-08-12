package clients_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/felipeosantos/curso-go/goexpert/observabilidade/internal/infra/webserver/clients"
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

func TestBuscaClimaCEPB_Success(t *testing.T) {

	mockResponse := `{
	  "city": "Guarulhos",
	  "temp_C": 17.8,
	  "temp_F": 64.04,
	  "temp_K": 290.95
	}`
	mockClient := &MockHTTPClientWeatherAPI{
		Do: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(mockResponse)),
			}, nil
		},
	}

	result, climacepBError, httpStatusCode, err := clients.BuscaClimaCEPB(context.Background(), mockClient, "07050000")
	assert.NoError(t, err)
	assert.Nil(t, climacepBError)
	assert.NotNil(t, result)
	assert.Equal(t, "Guarulhos", result.City)
	assert.Equal(t, 17.8, result.TempCelsius)
	assert.Equal(t, http.StatusOK, httpStatusCode)
}

func TestBuscaClimaCEPB_HTTPClientError(t *testing.T) {

	mockClient := &MockHTTPClientWeatherAPI{
		Do: func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("network error")
		},
	}

	result, climacepBError, httpStatusCode, err := clients.BuscaClimaCEPB(context.Background(), mockClient, "07050000")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, climacepBError)
	assert.Equal(t, http.StatusInternalServerError, httpStatusCode)
}

func TestBuscaClimaCEPB_StatusCodeError(t *testing.T) {

	mockResponse := `{
	  "mensagem": "can not find zipcode"
	}`
	mockClient := &MockHTTPClientWeatherAPI{
		Do: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(bytes.NewBufferString(mockResponse)),
			}, nil
		},
	}

	result, climacepBError, httpStatusCode, err := clients.BuscaClimaCEPB(context.Background(), mockClient, "0700000A")
	assert.NoError(t, err)
	assert.NotNil(t, climacepBError)
	assert.Nil(t, result)
	assert.Equal(t, "can not find zipcode", climacepBError.Mensagem)
	assert.Equal(t, http.StatusBadRequest, httpStatusCode)
}
