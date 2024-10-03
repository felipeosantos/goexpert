package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type EconomiaResp struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

type EconomiaErrorResp struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Cotacao struct {
	VlCambioDolar float64 `json:"vlCambioDolar"`
	gorm.Model    `json:"-"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	os.Remove("cotacao.db")
	mux := http.NewServeMux()
	mux.HandleFunc("GET /cotacao", BuscarCotacaoHandler)
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func BuscarCotacaoHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("Received request for " + r.URL.Path)

	cotacao, err := BuscarCotacao("USD-BRL")
	if err != nil {
		log.Printf("Error fetching cotacao: %v", err)
		writeJsonError(w, err, http.StatusInternalServerError)
		return
	}

	if err := GravarCotacao(cotacao); err != nil {
		log.Printf("Error saving cotacao: %v", err)
	}

	if err := json.NewEncoder(w).Encode(cotacao); err != nil {
		log.Printf("Error encoding response: %v", err)
		writeJsonError(w, err, http.StatusInternalServerError)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		log.Println("Cotacao fetched successfully")
		return
	}
}

func BuscarCotacao(moeda string) (*Cotacao, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/"+moeda, nil)
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

		var result EconomiaErrorResp
		err = json.Unmarshal(body, &result)
		if err != nil {
			log.Printf("Error parsing error response: %v", err)
			return nil, err
		}

		return nil, fmt.Errorf("Error response from API: %v - %v", result.Code, result.Message)

	}

	var result map[string]EconomiaResp
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("Error parsing response: %v", err)
		return nil, err
	}

	bidValue, err := strconv.ParseFloat(result[strings.Replace(moeda, "-", "", -1)].Bid, 64)
	if err != nil {
		log.Printf("Error parsing bid value: %v", err)
		return nil, err
	}

	return &Cotacao{
		VlCambioDolar: bidValue,
	}, nil
}

func GravarCotacao(cotacao *Cotacao) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	db, err := gorm.Open(sqlite.Open("cotacao.db"), &gorm.Config{})
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return err
	}

	// Migrate the schema
	if err := db.AutoMigrate(&Cotacao{}); err != nil {
		log.Printf("Error migrating schema: %v", err)
		return err
	}

	// Create
	if err := db.WithContext(ctx).Create(cotacao).Error; err != nil {
		log.Printf("Error creating cotacao: %v", err)
		return err
	}

	return nil
}

func writeJsonError(w http.ResponseWriter, err error, statusCode int) {
	errorResponse := ErrorResponse{Error: err.Error()}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		log.Printf("Error encoding error response: %v", err)
	}
}
