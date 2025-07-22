package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/felipeosantos/curso-go/goexpert/cloud-run/internal/entity"
)

func BuscaCEP(httpClient HTTPClient, cep string) (*entity.ViaCEP, error) {

	viacepURL := &url.URL{
		Scheme: "http",
		Host:   "viacep.com.br",
		Path:   fmt.Sprintf("/ws/%s/json/", url.PathEscape(cep)),
	}
	req, err := http.NewRequest("GET", viacepURL.String(), nil)
	if err != nil {
		fmt.Printf("Erro ao criar requisição para ViaCEP: %v\n", err)
		return nil, err
	}

	resp, err := httpClient.DoFunc(req)
	if err != nil {
		fmt.Printf("Erro ao fazer requisição para ViaCEP: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		fmt.Printf("ViaCEP retornou status code inválido: %d\n", resp.StatusCode)
		return nil, fmt.Errorf("erro ao buscar CEP - status code ViaCEP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Erro ao ler resposta da ViaCEP: %v\n", err)
		return nil, err
	}

	var viaCEP entity.ViaCEP
	err = json.Unmarshal(body, &viaCEP)
	if err != nil {
		fmt.Printf("Erro ao fazer parse da resposta da ViaCEP: %v\n", err)
		return nil, err
	}

	return &viaCEP, nil
}
