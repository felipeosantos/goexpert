/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

type Result struct {
	ID         int
	StatusCode int
	Duration   time.Duration
	Error      error
}

var (
	url         string
	requests    int
	concurrency int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "stress-test",
	Short: "Ferramenta de teste de carga para serviços web",
	Long: `Uma aplicação CLI escrita em Go para realizar testes de carga em serviços web.
	
Esta ferramenta permite configurar o número total de requisições e o nível de
concorrência, gerando relatórios detalhados sobre o desempenho do serviço testado.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Here you can add any pre-run checks or setup if needed
		// Validate requests and concurrency greater than zero
		if requests <= 0 {
			return errors.New("o número de requisições deve ser um inteiro positivo")
		}
		if concurrency <= 0 {
			return errors.New("o nível de concorrência deve ser um inteiro positivo")
		}
		// requests should be greater than or equal to concurrency
		if requests < concurrency {
			return errors.New("o número de requisições deve ser maior ou igual ao nível de concorrência")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		//url, _ := cmd.Flags().GetString("url")
		//requests, _ := cmd.Flags().GetInt("requests")
		//concurrency, _ := cmd.Flags().GetInt("concurrency")
		fmt.Printf("Iniciando teste de stress com os seguintes parâmetros:\n")
		fmt.Printf("URL: %s\n", url)
		fmt.Printf("Total de requisições: %d\n", requests)
		fmt.Printf("Nível de concorrência: %d\n", concurrency)
		fmt.Println("----------------------------------------")

		startTime := time.Now()
		results := executeTest(url, requests, concurrency)
		totalDuration := time.Since(startTime)

		generateReport(results, totalDuration, requests)
		fmt.Printf("Duração total do teste: %v\n", totalDuration)

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.stress-test.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	// Flags URL / requests / concurrency
	rootCmd.Flags().StringVarP(&url, "url", "u", "", "URL alvo para o teste de stress")
	rootCmd.Flags().IntVarP(&requests, "requests", "r", 0, "Número de requisições a serem realizadas")
	rootCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 0, "Número de requisições simultâneas")
	// Mark flags as required
	rootCmd.MarkFlagsRequiredTogether("url", "requests", "concurrency")
}

func worker(id int, url string, jobs <-chan int, results chan<- Result) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	for job := range jobs {
		//fmt.Printf("Worker %d processando job: %d\n", id, job)
		startTime := time.Now()
		resp, err := client.Get(url)
		duration := time.Since(startTime)

		result := Result{
			ID:       job,
			Duration: duration,
		}

		if err != nil {
			result.StatusCode = 0
			result.Error = err
		} else {
			result.StatusCode = resp.StatusCode
			resp.Body.Close()
		}

		results <- result
	}
}

func executeTest(url string, totalRequests, concurrencyLevel int) []Result {
	results := make([]Result, totalRequests)

	// Usar um canal não-bufferizado para jobs - aplica backpressure natural
	jobs := make(chan int)

	// Usar um canal com buffer moderado para resultados
	// Buffer com tamanho = concurrencyLevel * 2 é geralmente um bom equilíbrio
	resultsChan := make(chan Result, concurrencyLevel*2)

	// Iniciar workers
	for w := 0; w < concurrencyLevel; w++ {
		go worker(w, url, jobs, resultsChan)
	}

	// Iniciar uma goroutine para coletar resultados
	done := make(chan bool)
	go func() {
		for i := 0; i < totalRequests; i++ {
			result := <-resultsChan
			results[result.ID] = result
		}
		done <- true
	}()

	// Enviar jobs para os workers
	for j := 0; j < totalRequests; j++ {
		jobs <- j
	}
	close(jobs)

	// Esperar todos os resultados serem processados
	<-done
	close(resultsChan)

	return results
}

func generateReport(results []Result, totalDuration time.Duration, totalRequests int) {
	// Count status codes
	statusCodes := make(map[int]int)
	var totalSuccessful, totalFailed int
	var totalResponseTime time.Duration

	for _, result := range results {
		statusCodes[result.StatusCode]++
		if result.Error == nil {
			totalResponseTime += result.Duration
			if result.StatusCode == 200 {
				totalSuccessful++
			} else {
				totalFailed++
			}
		} else {
			totalFailed++
		}
	}

	// Print report
	fmt.Println("\nRelatório de Resultados do Teste")
	fmt.Println("----------------------------------------")
	fmt.Printf("Tempo total: %v\n", totalDuration)
	fmt.Printf("Total de requisições: %d\n", totalRequests)
	fmt.Printf("Requisições bem-sucedidas (HTTP 200): %d\n", totalSuccessful)
	fmt.Printf("Requisições com falha: %d\n", totalFailed)

	if totalSuccessful > 0 {
		fmt.Printf("Tempo médio de resposta: %v\n", totalResponseTime/time.Duration(totalSuccessful))
	}

	fmt.Println("\nDistribuição de Códigos de Status:")
	for code, count := range statusCodes {
		if code == 0 {
			fmt.Printf("  Erros de conexão: %d\n", count)
		} else {
			fmt.Printf("  HTTP %d: %d\n", code, count)
		}
	}

	fmt.Printf("\nRequisições por segundo: %.2f\n", float64(totalRequests)/totalDuration.Seconds())
}
