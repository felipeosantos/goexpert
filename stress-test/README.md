# Ferramenta de Teste de Carga para Serviços Web

Uma aplicação CLI escrita em Go para realizar testes de carga em serviços web.

## Funcionalidades

- Realizar testes de stress em serviços web
- Configurar número total de requisições
- Definir nível de concorrência
- Gerar relatórios detalhados incluindo:
  - Tempo total de execução
  - Total de requisições executadas
  - Número de requisições bem-sucedidas (HTTP 200)
  - Distribuição de códigos de status HTTP
  - Requisições por segundo

## Como Usar

### Execução local:

```bash
go run main.go --url=https://example.com --requests=1000 --concurrency=10
```

### Via Docker:

Construir a imagem:
```bash
docker build -t stress-test .
```

Executar o container:
```bash
docker run stress-test --url=http://example.com --requests=1000 --concurrency=10
```

## Parâmetros

- `--url`: URL do serviço a ser testado
- `--requests`: Número total de requisições a serem realizadas
- `--concurrency`: Número de chamadas simultâneas

## Exemplo de Saída

```
Iniciando teste de stress com os seguintes parâmetros:
URL: https://example.com
Total de requisições: 1000
Nível de concorrência: 10
----------------------------------------

Relatório de Resultados do Teste
----------------------------------------
Tempo total: 15.234s
Total de requisições: 1000
Requisições bem-sucedidas (HTTP 200): 982
Requisições com falha: 18
Tempo médio de resposta: 152.341ms

Distribuição de Códigos de Status:
  HTTP 200: 982
  HTTP 404: 12
  HTTP 500: 6

Requisições por segundo: 65.64
```