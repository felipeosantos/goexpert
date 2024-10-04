package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Cotacao struct {
	VlCambioDolar float64 `json:"vlCambioDolar"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	log.Println("Starting client...")
	cotacao, err := BuscarCotacao()
	if err != nil {
		log.Printf("Error fetching cotacao: %v", err)
		return
	}

	if err := GravarArquivo(cotacao); err != nil {
		log.Printf("Error saving cotacao: %v", err)
	}
}

func BuscarCotacao() (*Cotacao, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return nil, err
	}

	if resp.StatusCode > 299 {

		var result ErrorResponse
		err = json.Unmarshal(body, &result)
		if err != nil {
			log.Printf("Error parsing error response: %v", err)
			return nil, err
		}

		return nil, fmt.Errorf("Error response from API: %v", result.Error)

	}

	var result Cotacao
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("Error parsing response: %v", err)
		return nil, err
	}

	return &result, nil
}

func GravarArquivo(cotacao *Cotacao) error {

	file, err := os.Create("cotacao.txt")
	if err != nil {
		log.Printf("Error creating file: %v", err)
		return err
	}

	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Dólar: %v", cotacao.VlCambioDolar))
	if err != nil {
		log.Printf("Error writing to file: %v", err)
		return err
	}
	log.Println("File created successfully!")
	return nil
}
