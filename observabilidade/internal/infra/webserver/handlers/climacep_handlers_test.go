package handlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/felipeosantos/curso-go/goexpert/observabilidade/internal/entity"
	"github.com/felipeosantos/curso-go/goexpert/observabilidade/internal/infra/webserver/handlers"
	"github.com/stretchr/testify/assert"
)

// MockHTTPClientClimaCEPB implements HTTPClient for testing
type MockHTTPClientClimaCEPB struct {
	Do func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClientClimaCEPB) DoFunc(req *http.Request) (*http.Response, error) {
	return m.Do(req)
}

func TestBuscaClimaCEP_Success(t *testing.T) {

	mockResponseClimaCEPB := `{
	  "city": "Guarulhos",
	  "temp_C": 20.4,
	  "temp_F": 68.72,
	  "temp_K": 293.55
	}`

	mockClientClimaCEPB := &MockHTTPClientClimaCEPB{
		Do: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(mockResponseClimaCEPB)),
			}, nil
		},
	}

	handler := &handlers.ClimaCEPHandler{
		HTTPClient: mockClientClimaCEPB,
	}

	requestBody := map[string]string{
		"cep": "01001000",
	}
	bodyBytes, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Create a ResponseRecorder to record the response
	w := httptest.NewRecorder()

	handler.BuscaClimaCEP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp entity.ClimaCEPResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "Guarulhos", resp.City)
	assert.Equal(t, 20.4, resp.TempCelsius)
	assert.Equal(t, 68.72, resp.TempFahrenheit)
	assert.Equal(t, 293.55, resp.TempKelvin)
}

func TestBuscaClimaCEP_InvalidCEP(t *testing.T) {

	mockResponseClimaCEPB := `		{
	  "mensagem": "invalid zipcode"
	}`

	mockClientClimaCEPB := &MockHTTPClientClimaCEPB{
		Do: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusUnprocessableEntity,
				Body:       io.NopCloser(bytes.NewBufferString(mockResponseClimaCEPB)),
			}, nil
		},
	}

	handler := &handlers.ClimaCEPHandler{
		HTTPClient: mockClientClimaCEPB,
	}

	requestBody := map[string]string{
		"cep": "0700000A",
	}
	bodyBytes, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	// Create a ResponseRecorder to record the response
	w := httptest.NewRecorder()

	handler.BuscaClimaCEP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	var resp entity.ClimaCEPError
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "invalid zipcode", resp.Mensagem)
}

func TestBuscaClimaCEP_NotFoundCEP(t *testing.T) {

	mockResponseClimaCEPB := `{
	  "mensagem": "can not find zipcode"
	}`

	mockClientClimaCEPB := &MockHTTPClientClimaCEPB{
		Do: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(bytes.NewBufferString(mockResponseClimaCEPB)),
			}, nil
		},
	}

	handler := &handlers.ClimaCEPHandler{
		HTTPClient: mockClientClimaCEPB,
	}

	requestBody := map[string]string{
		"cep": "07000000",
	}
	bodyBytes, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	// Create a ResponseRecorder to record the response
	w := httptest.NewRecorder()

	handler.BuscaClimaCEP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var resp entity.ClimaCEPError
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "can not find zipcode", resp.Mensagem)
}

func TestBuscaClimaCEP_UnprocessableEntity(t *testing.T) {

	mockResponseClimaCEPB := `{
	  "mensagem": "invalid zipcode"
	}`

	mockClientClimaCEPB := &MockHTTPClientClimaCEPB{
		Do: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusUnprocessableEntity,
				Body:       io.NopCloser(bytes.NewBufferString(mockResponseClimaCEPB)),
			}, nil
		},
	}

	handler := &handlers.ClimaCEPHandler{
		HTTPClient: mockClientClimaCEPB,
	}

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	// Create a ResponseRecorder to record the response
	w := httptest.NewRecorder()

	handler.BuscaClimaCEP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	var resp entity.ClimaCEPError
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "invalid zipcode", resp.Mensagem)
}
