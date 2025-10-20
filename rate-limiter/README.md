# Limitador de Requisições em Go

Um middleware configurável de limitação de requisições para aplicações web em Go que pode restringir o acesso com base no endereço IP ou token de API.

## Funcionalidades

- Limitação de requisições baseada em endereço IP
- Limitação de requisições baseada em token de API
- Limites de taxa e durações de bloqueio configuráveis
- Suporte para armazenamento em Redis ou em memória
- Fácil integração com o roteador Chi
- Configuração através de variáveis de ambiente ou arquivo .env

## Configuração

O limitador de requisições pode ser configurado usando variáveis de ambiente ou um arquivo `.env` com estrutura hierárquica:

Obs: As variáveis de ambiente têm precedência sobre as configurações do arquivo `.env`. Para as variáveis de ambiente substituir `.` por `_`.

### Configuração do Servidor

| Variável | Descrição | Padrão |
|----------|-------------|---------|
| SERVER_PORT | Porta para o servidor | 8080 |

### Configuração de Armazenamento

| Variável | Descrição | Padrão |
|----------|-------------|---------|
| STORAGE_TYPE | Backend de armazenamento (redis ou memory) | memory |
| STORAGE_REDIS_URL | URL para conexão com Redis | redis://localhost:6379/0 |
| STORAGE_MEMORY_SIZE | Número máximo de entradas para armazenamento em memória | 1000 |

### Configuração de Limitação por IP

| Variável | Descrição | Padrão |
|----------|-------------|---------|
| IP_RATE_LIMIT | Requisições permitidas por IP na janela de tempo | 1 |
| IP_RATE_WINDOW | Janela de tempo para limitação de IP | 1s |
| IP_BLOCK_DURATION | Quanto tempo bloquear o IP após exceder o limite | 10s |

### Configuração Específica por Token

Configure limites de requisição personalizados para tokens específicos usando o formato hierárquico:

```
TOKEN.[nome_token].RATE_LIMIT=[número]
TOKEN.[nome_token].RATE_WINDOW=[duração]
TOKEN.[nome_token].BLOCK_DURATION=[duração]
```

#### Exemplos de Configuração de Token

Obs: Os tokens devem estar em caixa baixa no envio das requisições.

```env
# Configuração do token "acb"
TOKEN.ACB.RATE_LIMIT=2
TOKEN.ACB.RATE_WINDOW=1s
TOKEN.ACB.BLOCK_DURATION=20s

# Configuração do token "devwas"
TOKEN.DEVWAS.RATE_LIMIT=50
TOKEN.DEVWAS.RATE_WINDOW=1s
TOKEN.DEVWAS.BLOCK_DURATION=1m

# Configuração do token "asdqwed"
TOKEN.ASDQWED.RATE_LIMIT=20
TOKEN.ASDQWED.RATE_WINDOW=1s
TOKEN.ASDQWED.BLOCK_DURATION=1m
```

### Formato de Duração

Os valores de duração podem ser especificados usando o formato de duração do Go:
- `s` para segundos (ex: `10s`)
- `m` para minutos (ex: `1m`)
- `h` para horas (ex: `1h`)
- Formato combinado (ex: `1m30s`)

## Executando a Aplicação

### Desenvolvimento Local

```bash
# Executar a aplicação
go run cmd/server/main.go
```

### Docker Compose

```bash
# Iniciar com Redis e os contêineres da aplicação
docker-compose up -d
```

## Testes

Você pode testar o limitador de requisições usando o arquivo HTTP no diretório `test`:

```
test/api.http
```

Ou usando curl:

```bash
# Testar limitação baseada em IP
curl http://localhost:8080/

# Testar limitação baseada em token
curl -H "API_KEY: acb" http://localhost:8080/
curl -H "API_KEY: devwas" http://localhost:8080/
curl -H "API_KEY: asdqwed" http://localhost:8080/

# Disparar limitação de requisições (enviar múltiplas requisições rapidamente)
for i in {1..5}; do curl http://localhost:8080/; done
for i in {1..3}; do curl -H "API_KEY: acb" http://localhost:8080/; done
```

## Arquitetura

O limitador de requisições segue um design modular:

- **Padrão Strategy**: A interface de armazenamento pode ser implementada por diferentes backends (Redis, em memória)
- **Padrão Middleware**: O limitador de requisições pode ser injetado na cadeia de handlers HTTP
- **Configuração**: Variáveis de ambiente hierárquicas com notação de ponto
- **Padrão Factory**: Cria armazenamento com base na configuração

## Estrutura do Pacote

```
├── cmd/
│   └── server/          # Ponto de entrada da aplicação
├── config/              # Gerenciamento de configuração
├── internal/
│   ├── limiter/         # Lógica central de limitação de requisições
│   ├── middleware/      # Implementação de middleware HTTP
│   └── storage/         # Implementações de armazenamento (Redis, em memória)
├── test/                # Arquivos de teste e exemplos de API
├── .env                 # Configuração de ambiente com estrutura hierárquica
└── docker-compose.yml   # Composição Docker
```

## Exemplo de Arquivo .env

```env
# Configuração do servidor
SERVER_PORT=8080

# Tipo de armazenamento (redis ou memory)
STORAGE_TYPE=memory
# Configuração de armazenamento Redis
STORAGE.REDIS.URL=redis://localhost:6379/0

# Configuração de armazenamento em memória
# STORAGE.MEMORY.SIZE=1000

# Configuração padrão de limitação de IP
IP.RATE_LIMIT=1
IP.RATE_WINDOW=1s
IP.BLOCK_DURATION=10s

# Configurações de limitação por token
TOKEN.ACB.RATE_LIMIT=2
TOKEN.ACB.RATE_WINDOW=1s
TOKEN.ACB.BLOCK_DURATION=20s

TOKEN.DEVWAS.RATE_LIMIT=50
TOKEN.DEVWAS.RATE_WINDOW=1s
TOKEN.DEVWAS.BLOCK_DURATION=1m

TOKEN.ASDQWED.RATE_LIMIT=20
TOKEN.ASDQWED.RATE_WINDOW=1s
TOKEN.ASDQWED.BLOCK_DURATION=1m
```
