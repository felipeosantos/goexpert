package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/felipeosantos/curso-go/goexpert/observabilidade/internal/entity"
)

func BuscaClimaCEPB(ctx context.Context, httpClient HTTPClient, cep string) (*entity.ClimaCEPResponse, *entity.ClimaCEPError, int, error) {

	baseURL := os.Getenv("CLIMA_CEP_B_BASE_URL")
	params := url.Values{}
	params.Set("cep", cep)

	url := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Printf("Erro ao criar requisição para WeatherAPI: %v\n", err)
		return nil, nil, http.StatusInternalServerError, err
	}

	resp, err := httpClient.DoFunc(req)
	if err != nil {
		fmt.Printf("Erro ao fazer requisição para WeatherAPI: %v\n", err)
		return nil, nil, http.StatusInternalServerError, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Erro ao ler resposta da WeatherAPI: %v\n", err)
		return nil, nil, http.StatusInternalServerError, err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		var climacepBError entity.ClimaCEPError
		err = json.Unmarshal(body, &climacepBError)
		if err != nil {
			fmt.Printf("Erro ao fazer parse da resposta da WeatherAPI: %v\n", err)
			return nil, nil, http.StatusInternalServerError, err
		}

		return nil, &climacepBError, resp.StatusCode, nil
	} else {
		var climacepBResponse entity.ClimaCEPResponse

		err = json.Unmarshal(body, &climacepBResponse)
		if err != nil {
			fmt.Printf("Erro ao fazer parse da resposta da WeatherAPI: %v\n", err)
			return nil, nil, http.StatusInternalServerError, err
		}

		return &climacepBResponse, nil, resp.StatusCode, nil

	}

}
