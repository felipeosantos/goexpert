package entity

// Simple configuration struct
//type configEnv struct {
//	WeatherApiKey string `env:"WEATHER_API_KEY"` // Weather API key
//}

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
type ClimaResult struct {
	TempCelsius    float64 `json:"temp_C"`
	TempFahrenheit float64 `json:"temp_F"`
	TempKelvin     float64 `json:"temp_K"`
}

type ErrorResponse struct {
	Mensagem string `json:"mensagem"`
}

type CurrentWeatherResponse struct {
	Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Lat            float64 `json:"lat"`
		Lon            float64 `json:"lon"`
		TzID           string  `json:"tz_id"`
		LocaltimeEpoch int     `json:"localtime_epoch"`
		Localtime      string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		LastUpdatedEpoch int     `json:"last_updated_epoch"`
		LastUpdated      string  `json:"last_updated"`
		TempC            float64 `json:"temp_c"`
		TempF            float64 `json:"temp_f"`
		IsDay            int     `json:"is_day"`
		Condition        struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
			Code int    `json:"code"`
		} `json:"condition"`
		WindMph    float64 `json:"wind_mph"`
		WindKph    float64 `json:"wind_kph"`
		WindDegree float64 `json:"wind_degree"`
		WindDir    string  `json:"wind_dir"`
		PressureMb float64 `json:"pressure_mb"`
		PressureIn float64 `json:"pressure_in"`
		PrecipMm   float64 `json:"precip_mm"`
		PrecipIn   float64 `json:"precip_in"`
		Humidity   float64 `json:"humidity"`
		Cloud      float64 `json:"cloud"`
		FeelslikeC float64 `json:"feelslike_c"`
		FeelslikeF float64 `json:"feelslike_f"`
		VisKm      float64 `json:"vis_km"`
		VisMiles   float64 `json:"vis_miles"`
		Uv         float64 `json:"uv"`
		GustMph    float64 `json:"gust_mph"`
		GustKph    float64 `json:"gust_kph"`
		AirQuality struct {
			Co           float64 `json:"co"`
			No2          float64 `json:"no2"`
			O3           float64 `json:"o3"`
			So2          float64 `json:"so2"`
			Pm25         float64 `json:"pm2_5"`
			Pm10         float64 `json:"pm10"`
			UsEpaIndex   int     `json:"us-epa-index"`
			GbDefraIndex int     `json:"gb-defra-index"`
		} `json:"air_quality"`
	} `json:"current"`
}

type WeatherErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	HttpStatusCode int `json:"http_status_code"` // Adicionando o campo para armazenar o status code HTTP
}
