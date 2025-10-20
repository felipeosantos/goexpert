# Rate Limiter in Go

A configurable rate limiter middleware for Go web applications that can restrict access based on IP address or API token.

## Features

- Limit requests based on IP address
- Limit requests based on API token
- Configurable rate limits and block durations
- Support for Redis or in-memory storage
- Easy integration with Chi router
- Configuration via environment variables or .env file

## Configuration

The rate limiter can be configured using environment variables or a `.env` file:

| Variable | Description | Default |
|----------|-------------|---------|
| SERVER_PORT | Port for the server | 8080 |
| REDIS_URL | URL for Redis connection | redis://localhost:6379/0 |
| STORAGE_TYPE | Storage backend (redis or memory) | redis |
| IP_RATE_LIMIT | Requests per second allowed per IP | 10 |
| TOKEN_RATE_LIMIT | Requests per second allowed per token | 100 |
| BLOCK_DURATION | How long to block after limit exceeded | 5m |
| RATE_WINDOW | Time window for rate limiting | 1s |
| TOKEN_CONFIG_[token_name] | Token-specific rate limit and block duration | N/A |

### Token-specific Configuration

You can configure custom rate limits and block durations for specific tokens using the following format:

```
TOKEN_CONFIG_[token_name]=[rate_limit],[block_duration]
```

For example:
```
TOKEN_CONFIG_premium=200,10m  # Premium token: 200 req/s, 10 minute block
TOKEN_CONFIG_basic=50,5m      # Basic token: 50 req/s, 5 minute block
TOKEN_CONFIG_test=20,1m       # Test token: 20 req/s, 1 minute block
```

## Running with Docker Compose

To run the application with Docker Compose:

```bash
docker-compose up -d
```

This will start both Redis and the application containers.

## Testing

You can test the rate limiter using the HTTP file in the `test` directory:

```
test/api.http
```

Or using curl:

```bash
# Test regular request (limited by IP)
curl http://localhost:8080/

# Test with API key (limited by token)
curl -H "API_KEY: test-token-123" http://localhost:8080/

# Send multiple requests in quick succession to trigger rate limiting
for i in {1..20}; do curl -H "API_KEY: test-token-123" http://localhost:8080/; done
```

## Architecture

The rate limiter follows a modular design:

- **Strategy Pattern**: Storage interface can be implemented by different backends
- **Middleware Pattern**: Rate limiter can be injected into the HTTP handler chain
- **Configuration**: Environment variables allow flexible configuration

## Package Structure

- `/cmd/server`: Main application entry point
- `/config`: Configuration loading from environment variables
- `/internal/limiter`: Core rate limiting logic
- `/internal/middleware`: HTTP middleware for rate limiting
- `/internal/storage`: Storage implementations (Redis, in-memory)
- `/test`: Test files
