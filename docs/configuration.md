# Configuration

Runtime configuration is loaded from `configs/config.yml` and can be overridden by environment variables.

## YAML config

```yaml
grpc:
  port: 50051

http:
  port: 8080

database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  name: go_api
  ssl_mode: disable

logger:
  level: info
  format: json
```

## Override config path

```bash
GOAPI_CONFIG=/etc/go-api/config.yml make run
```

## Environment variable override

```bash
GOAPI_GRPC_PORT=50051
GOAPI_HTTP_PORT=8080
GOAPI_DB_HOST=localhost
GOAPI_DB_PORT=5432
GOAPI_DB_USER=postgres
GOAPI_DB_PASSWORD=postgres
GOAPI_DB_NAME=go_api
GOAPI_DB_SSL_MODE=disable
GOAPI_LOGGER_LEVEL=info
GOAPI_LOGGER_FORMAT=json
```

Fallback aliases are also supported:
`DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSL_MODE`, `GRPC_PORT`, `HTTP_PORT`, `LOGGER_LEVEL`, `LOGGER_FORMAT`.

## Notes

- Environment variables have higher priority than YAML values.
- Service can still start with env-only config when config file is missing.

## Kubernetes secret example

```yaml
env:
  - name: GOAPI_DB_HOST
    value: postgres.default.svc.cluster.local
  - name: GOAPI_DB_PORT
    value: '5432'
  - name: GOAPI_DB_USER
    valueFrom:
      secretKeyRef:
        name: go-api-secrets
        key: db_user
  - name: GOAPI_DB_PASSWORD
    valueFrom:
      secretKeyRef:
        name: go-api-secrets
        key: db_password
  - name: GOAPI_DB_NAME
    value: go_api
```
