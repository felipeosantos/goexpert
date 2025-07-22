# ClimaCEP API

Uma API REST em Go que consulta informações de clima baseadas em um CEP brasileiro. A aplicação integra as APIs do ViaCEP para validação de endereços e WeatherAPI para obtenção de dados meteorológicos.

## 🌐 Demo Online

A aplicação está disponível online no Google Cloud Run:

**🔗 https://goclimacep-leizos6cda-uc.a.run.app**

Teste diretamente no seu navegador:
- [CEP São Paulo - Centro](https://goclimacep-leizos6cda-uc.a.run.app/?cep=01001000)
- [CEP Campo Grande - MS](https://goclimacep-leizos6cda-uc.a.run.app/?cep=79002000)

## 📋 Funcionalidades

- Consulta de CEP através da API ViaCEP
- Obtenção de dados meteorológicos através da WeatherAPI
- Conversão automática de temperaturas (Celsius, Fahrenheit, Kelvin)
- Validação de formato de CEP
- Tratamento de erros robusto
- Containerização com Docker

## 🏗️ Arquitetura do Projeto

```
cloud-run/
├── cmd/server/                 # Ponto de entrada da aplicação
│   └── main.go                # Configuração do servidor HTTP
├── internal/
│   ├── entity/                # Definições de estruturas de dados
│   │   └── entity.go          # Entidades para APIs externas
│   └── infra/webserver/       # Infraestrutura web
│       ├── clients/           # Clientes HTTP para APIs externas
│       │   ├── http_client.go # Cliente HTTP customizado
│       │   ├── viacep.go      # Cliente para API ViaCEP
│       │   └── weatherapi.go  # Cliente para WeatherAPI
│       └── handlers/          # Handlers HTTP
│           └── climacep_handlers.go # Lógica principal da API
├── test/
│   └── api.http              # Testes de API REST
├── docker-compose.yaml       # Configuração para desenvolvimento
├── Dockerfile               # Docker para desenvolvimento
├── Dockerfile.prod          # Docker para produção
├── go.mod                   # Dependências do projeto
└── README.md               # Documentação
```

## 🚀 Como Executar

### Pré-requisitos

- Go 1.21.13 ou superior
- Docker e Docker Compose (opcional)
- Chave de API do WeatherAPI (gratuita em https://www.weatherapi.com/)

### 1. Configuração da API Key

Crie um arquivo `.env` em `cmd/server/.env` com o seguinte conteúdo:

```env
WEATHER_API_KEY=sua_chave_aqui
```

### 2. Execução Local (Go)

```bash
# Instalar dependências
go mod tidy

# Executar a aplicação
cd cmd/server
go run main.go
```

A API estará disponível em `http://localhost:8080`

### 3. Execução com Docker Compose (Desenvolvimento)

```bash
docker-compose up -d --build
```

Este comando:
- Constrói a imagem usando o `Dockerfile` (modo desenvolvimento)
- Inicia o container com volume montado para desenvolvimento
- Expõe a aplicação na porta 8080

### 4. Build para Produção

```bash
# Construir imagem de produção
docker build -t <your_dockerhub_username>/<your_image_name> -f Dockerfile.prod .

# Executar container de produção
docker run --rm --env WEATHER_API_KEY=<your_api_key> -p 8080:8080 <your_dockerhub_username>/<your_image_name>:<your_image_tag>
```

## 📚 Como Usar a API

### Endpoint Principal

**Local:**
```
GET http://localhost:8080/?cep={cep}
```

**Produção (Google Cloud Run):**
```
GET https://goclimacep-leizos6cda-uc.a.run.app/?cep={cep}
```

### Exemplos de Uso

#### Sucesso (CEP válido):

**Local:**
```bash
curl "http://localhost:8080/?cep=01001000"
```

**Produção:**
```bash
curl "https://goclimacep-leizos6cda-uc.a.run.app/?cep=01001000"
```

**Resposta:**
```json
{
  "temp_C": 25.0,
  "temp_F": 77.0,
  "temp_K": 298.15
}
```

#### Erro - CEP inválido:

**Local:**
```bash
curl "http://localhost:8080/?cep=123"
```

**Produção:**
```bash
curl "https://goclimacep-leizos6cda-uc.a.run.app/?cep=123"
```

**Resposta (422):**
```json
{
  "mensagem": "invalid zipcode"
}
```

#### Erro - CEP não encontrado:

**Local:**
```bash
curl "http://localhost:8080/?cep=99999999"
```

**Produção:**
```bash
curl "https://goclimacep-leizos6cda-uc.a.run.app/?cep=99999999"
```

**Resposta (404):**
```json
{
  "mensagem": "can not find zipcode"
}
```

## 🧪 Testes

Execute os testes automatizados:

```bash
go test ./...
```

Para testes manuais, use o arquivo `test/api.http` com extensões como REST Client no VS Code.

## 🐳 Docker

### Dockerfile (Desenvolvimento)
- Baseado em `golang:latest`
- Inclui volume para desenvolvimento
- Permite modificações em tempo real

### Dockerfile.prod (Produção)
- Multi-stage build para otimização
- Imagem final baseada em `scratch` (mínima)
- Tamanho reduzido e maior segurança
- Variável de ambiente para API key

## 🔧 Dependências

- **github.com/joho/godotenv**: Carregamento de variáveis de ambiente
- **github.com/stretchr/testify**: Framework de testes

## 📝 Estrutura das APIs Integradas

### ViaCEP API
- **Endpoint**: `https://viacep.com.br/ws/{cep}/json/`
- **Função**: Validação e obtenção de dados de endereço

### WeatherAPI
- **Endpoint**: `http://api.weatherapi.com/v1/current.json`
- **Função**: Obtenção de dados meteorológicos atuais
- **Requer**: API Key gratuita

## 🌐 Deploy

Este projeto está preparado para deploy em plataformas como:
- Google Cloud Run
- AWS ECS/Fargate
- Azure Container Instances
- Heroku
- Qualquer orquestrador Docker/Kubernetes

Para Cloud Run especificamente, a aplicação escuta na porta definida pela variável de ambiente `PORT` (padrão: 8080).