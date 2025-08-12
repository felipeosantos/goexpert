package clients_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/felipeosantos/curso-go/goexpert/cloud-run/internal/infra/webserver/clients"
	"github.com/stretchr/testify/assert"
)

// MockHTTPClientViaCEP implements HTTPClient for testing
type MockHTTPClientViaCEP struct {
	Do func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClientViaCEP) DoFunc(req *http.Request) (*http.Response, error) {
	return m.Do(req)
}

func TestBuscaCEP_Success(t *testing.T) {
	mockResponse := `{
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

	mockClient := &MockHTTPClientViaCEP{
		Do: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(mockResponse)),
			}, nil
		},
	}

	result, err := clients.BuscaCEP(context.Background(), mockClient, "01001000")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "01001-000", result.Cep)
	assert.Equal(t, "SP", result.Uf)
}

func TestBuscaCEP_HTTPClientError(t *testing.T) {
	mockClient := &MockHTTPClientViaCEP{
		Do: func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("network error")
		},
	}

	result, err := clients.BuscaCEP(context.Background(), mockClient, "01001000")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestBuscaCEP_StatusCodeError(t *testing.T) {
	mockClient := &MockHTTPClientViaCEP{
		Do: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(bytes.NewBufferString("")),
			}, nil
		},
	}

	result, err := clients.BuscaCEP(context.Background(), mockClient, "0000000")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestBuscaCEP_NotFoundError(t *testing.T) {
	mockResponse := `{
	  "erro": "true"
	}`

	mockClient := &MockHTTPClientViaCEP{
		Do: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(mockResponse)),
			}, nil
		},
	}

	result, err := clients.BuscaCEP(context.Background(), mockClient, "01000000")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "true", result.Erro)
}
