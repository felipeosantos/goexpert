package entity

// Estrutura de request de entrada para buscar clima por CEP
type ClimaCEPRequest struct {
	Cep string `json:"cep"`
}

// Estrutura para o resultado unificado
type ClimaCEPResponse struct {
	City           string  `json:"city"`
	TempCelsius    float64 `json:"temp_C"`
	TempFahrenheit float64 `json:"temp_F"`
	TempKelvin     float64 `json:"temp_K"`
}

type ClimaCEPError struct {
	Mensagem string `json:"mensagem"`
}
