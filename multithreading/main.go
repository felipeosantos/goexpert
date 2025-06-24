package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Estrutura para a resposta da API BrasilAPI
// Exemplo de resposta de sucesso da BrasilAPI, status code 200:
// {
//   "cep": "01001000",
//   "state": "SP",
//   "city": "São Paulo",
//   "neighborhood": "Sé",
//   "street": "Praça da Sé",
//   "service": "open-cep"
// }

type BrasilAPICEP struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

// Estrutura para a resposta de erro da API BrasilAPI
// Exemplo de resposta de erro da BrasilAPI, status code diferente de 200:
//
//	{
//	 "message": "Todos os serviços de CEP retornaram erro.",
//	 "type": "service_error",
//	 "name": "CepPromiseError",
//	 "errors": [
//	   {
//	     "name": "ServiceError",
//	     "message": "A autenticacao de null falhou!",
//	     "service": "correios"
//	   },
//	   {
//	     "name": "ServiceError",
//	     "message": "Cannot read properties of undefined (reading 'replace')",
//	     "service": "viacep"
//	   },
//	   {
//	     "name": "ServiceError",
//	     "message": "Erro ao se conectar com o serviço WideNet.",
//	     "service": "widenet"
//	   },
//	   {
//	     "name": "ServiceError",
//	     "message": "CEP não encontrado na base dos Correios.",
//	     "service": "correios-alt"
//	   }
//	 ]
//	}
type BrasilAPIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Errors  []struct {
		Name    string `json:"name"`
		Message string `json:"message"`
		Service string `json:"service"`
	} `json:"errors"`
}

// Estrutura para a resposta da API ViaCEP
// Exemplo de resposta de sucesso da ViaCEP, status code 200:
//
//	{
//	  "cep": "01001-000",
//	  "logradouro": "Praça da Sé",
//	  "complemento": "lado ímpar",
//	  "unidade": "",
//	  "bairro": "Sé",
//	  "localidade": "São Paulo",
//	  "uf": "SP",
//	  "estado": "São Paulo",
//	  "regiao": "Sudeste",
//	  "ibge": "3550308",
//	  "gia": "1004",
//	  "ddd": "11",
//	  "siafi": "7107"
//	}
//
// Exemplo de resposta de erro da ViaCEP, status code 200 (mas com campo "erro" verdadeiro):
//
//	{
//	  "erro": "true"
//	}
type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
	Erro        string `json:"erro"`
}

// Estrutura para o resultado unificado
type CEPResult struct {
	CEP          string
	Rua          string
	Bairro       string
	Cidade       string
	Estado       string
	APISource    string
	ResponseTime time.Duration
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Por favor, forneça um CEP como argumento. Exemplo: go run main.go 01001000")
		return
	}

	cep := os.Args[1]

	// Validando formato do CEP (formato mais completo)
	isValid := true
	if len(cep) != 8 {
		isValid = false
	}

	// Validando se foi informado numeros
	if _, err := strconv.Atoi(cep); err != nil {
		isValid = false
	}

	if !isValid {
		fmt.Println("Erro: CEP deve ter 8 dígitos numéricos, sem traço ou espaços.")
		return
	}

	fmt.Printf("Buscando informações para o CEP: %s\n\n", cep)

	// Criando um contexto com timeout de 1 segundo
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Canal para receber os resultados
	resultChan := make(chan CEPResult)

	// Consulta às duas APIs em paralelo
	go fetchBrasilAPI(ctx, cep, resultChan)
	go fetchViaCEP(ctx, cep, resultChan)

	// Esperando pelo resultado mais rápido
	select {
	case result := <-resultChan:
		printResult(result)
		return
	case <-ctx.Done():
		fmt.Println("Erro: Timeout de 1 segundo excedido. Nenhuma API respondeu a tempo.")
		return
	}
}

func fetchBrasilAPI(ctx context.Context, cep string, resultChan chan CEPResult) {
	url := "https://brasilapi.com.br/api/cep/v1/" + cep
	startTime := time.Now()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Printf("Erro ao criar requisição para BrasilAPI: %v\n", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Erro ao fazer requisição para BrasilAPI: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		fmt.Printf("BrasilAPI retornou status code inválido: %d\n", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Erro ao ler resposta da BrasilAPI: %v\n", err)
		return
	}

	var brasilAPICEP BrasilAPICEP
	err = json.Unmarshal(body, &brasilAPICEP)
	if err != nil {
		fmt.Printf("Erro ao fazer parse da resposta da BrasilAPI: %v\n", err)
		return
	}

	elapsedTime := time.Since(startTime)

	resultChan <- CEPResult{
		CEP:          brasilAPICEP.Cep,
		Rua:          brasilAPICEP.Street,
		Bairro:       brasilAPICEP.Neighborhood,
		Cidade:       brasilAPICEP.City,
		Estado:       brasilAPICEP.State,
		APISource:    "BrasilAPI",
		ResponseTime: elapsedTime,
	}
}

func fetchViaCEP(ctx context.Context, cep string, resultChan chan CEPResult) {
	url := "http://viacep.com.br/ws/" + cep + "/json/"
	startTime := time.Now()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Printf("Erro ao criar requisição para ViaCEP: %v\n", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Erro ao fazer requisição para ViaCEP: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		fmt.Printf("ViaCEP retornou status code inválido: %d\n", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Erro ao ler resposta da ViaCEP: %v\n", err)
		return
	}

	var viaCEP ViaCEP
	err = json.Unmarshal(body, &viaCEP)
	if err != nil {
		fmt.Printf("Erro ao fazer parse da resposta da ViaCEP: %v\n", err)
		return
	}

	// Verificar se a resposta contém um erro da API (ViaCEP retorna "erro": true para CEPs inválidos)
	if viaCEP.Erro == "true" {
		fmt.Println("ViaCEP retornou erro para o CEP informado")
		return
	}

	elapsedTime := time.Since(startTime)

	resultChan <- CEPResult{
		CEP:          viaCEP.Cep,
		Rua:          viaCEP.Logradouro,
		Bairro:       viaCEP.Bairro,
		Cidade:       viaCEP.Localidade,
		Estado:       viaCEP.Uf,
		APISource:    "ViaCEP",
		ResponseTime: elapsedTime,
	}
}

func printResult(result CEPResult) {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✅ Resultado recebido da API mais rápida")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("API: %s (respondeu em %v)\n", result.APISource, result.ResponseTime)
	fmt.Printf("CEP: %s\n", result.CEP)
	fmt.Printf("Rua: %s\n", result.Rua)
	fmt.Printf("Bairro: %s\n", result.Bairro)
	fmt.Printf("Cidade: %s\n", result.Cidade)
	fmt.Printf("Estado: %s\n", result.Estado)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
