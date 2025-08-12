# ClimaCEP API com Observabilidade

Uma API REST em Go que consulta informações de clima baseadas em um CEP brasileiro. Esta versão do projeto implementa observabilidade através de OpenTelemetry e Zipkin.

## 🏗️ Arquitetura do Projeto

```
observabilidade/
├── cmd/server/                 # Ponto de entrada da aplicação
│   └── main.go                # Configuração do servidor HTTP com OpenTelemetry
├── internal/
│   ├── entity/                # Definições de estruturas de dados
│   │   └── entity.go          # Entidades para APIs externas
│   └── infra/webserver/       # Infraestrutura web
│       ├── clients/           # Clientes HTTP para APIs externas
│       │   ├── http_client.go # Cliente HTTP customizado
│       │   └── climacepB.go   # Cliente para API ClimaCEP B
│       └── handlers/          # Handlers HTTP
│           └── climacep_handlers.go # Lógica principal da API
├── test/
│   └── api.http              # Testes de API REST
├── .docker/                  # Arquivos de configuração Docker
│   └── otel-collector.yaml   # Configuração do OpenTelemetry Collector
├── docker-compose.yaml       # Configuração para execução completa do sistema
├── Dockerfile               # Docker para desenvolvimento
├── Dockerfile.prod          # Docker para produção
├── go.mod                   # Dependências do projeto
└── README.md               # Documentação
```

## 📋 Sobre a Observabilidade

Este projeto implementa observabilidade completa utilizando:

- **OpenTelemetry**: Para instrumentação e coleta de telemetria
- **OpenTelemetry Collector**: Para processamento e exportação de dados
- **Zipkin**: Para visualização e análise de traces

A arquitetura consiste em dois serviços Go (serviceA e serviceB) que se comunicam entre si, com toda comunicação sendo instrumentada para rastreamento.

## 🚀 Como Executar

### Pré-requisitos

- Git
- Docker e Docker Compose

### Instruções

1. Fazer checkout do repositório:
```bash
git clone https://github.com/FelipeOSantos/goexpert.git
```

2. Navegar até a pasta do projeto de observabilidade:
```bash
cd goexpert/observabilidade
```

3. Executar os serviços usando Docker Compose:
```bash
docker-compose up -d --build
```

4. Executar chamadas usando o arquivo de teste:
   - Use o arquivo `test/api.http` se tiver a extensão REST Client no VS Code
   - Ou execute via curl:
     ```bash
     curl -X POST -H "Content-Type: application/json" -d '{"cep":"01001000"}' http://localhost:8081/
     ```

5. Abrir o navegador no Zipkin para visualizar os traces:
```
http://localhost:9411/
```
   - Na interface do Zipkin, clique em "Run Query" para ver os traces recentes
   - Explore os diferentes traces para visualizar a comunicação entre os serviços

6. Para parar os serviços:
```bash
docker-compose down
```

## 📚 Estrutura dos Serviços

- **serviceA (porta 8081)**: Serviço principal que recebe requisições e chama o serviceB
- **serviceB (porta 8080)**: Serviço que consulta APIs externas (ViaCEP e WeatherAPI)
- **OpenTelemetry Collector**: Recebe, processa e encaminha os dados de telemetria
- **Zipkin**: Interface web para visualização de traces

## 📊 Visualizando Traces

Após fazer algumas requisições à API, acesse o Zipkin em http://localhost:9411/ e:

1. Deixe os filtros padrão e clique em "Run Query"
2. Você verá os traces das requisições recentes
3. Clique em um trace para expandir e ver detalhes
4. Observe a cascata de chamadas entre os serviços
5. Verifique o tempo de resposta em cada etapa do processamento

## 🔍 Exemplo de Uso e Análise

1. Execute uma chamada à API:
```bash
curl -X POST -H "Content-Type: application/json" -d '{"cep":"01001000"}' http://localhost:8081/
```

2. No Zipkin, localize o trace correspondente
3. Observe como a requisição flui do serviceA para o serviceB
4. Analise os tempos de resposta e possíveis gargalos
5. Identifique os componentes da chamada e suas dependências

Esta implementação de observabilidade ajuda a entender o comportamento do sistema, diagnosticar problemas e otimizar o desempenho da aplicação.

