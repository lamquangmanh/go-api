# Quickstart

## Prerequisites

- Go 1.25+
- PostgreSQL
- make
- protoc
- sqlc

## Run locally

```bash
# from go-api/
make proto-tools
make proto
make sqlc
make tidy

# apply database migrations
make migrate-up DB_URL='postgresql://postgres:postgres@localhost:5432/go_api?sslmode=disable'

# run gRPC + HTTP health server
make run
```

## Build

```bash
make build
```

Output binary: `bin/grpc`

## Health check

```bash
curl -i http://localhost:8080/livez
curl -i http://localhost:8080/readyz
curl -i http://localhost:8080/healthz
```

Expected:

- `/livez` => `live`
- `/readyz` => `ok`
- `/healthz` => `ok`
