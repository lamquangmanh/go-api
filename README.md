# go-api

gRPC user/role management service built with Go, PostgreSQL, sqlc, and Protocol Buffers.

## Core Features

- gRPC `UserService` and `RoleService`
- gRPC health service + reflection enabled
- HTTP health endpoints (`/livez`, `/readyz`, `/healthz`)
- PostgreSQL persistence via `pgx` + `sqlc`
- Versioned DB migrations with `golang-migrate`

## Quick Start

```bash
# from go-api/
make proto-tools
make proto
make sqlc
make tidy
make migrate-up DB_URL='postgresql://postgres:postgres@localhost:5432/go_api?sslmode=disable'
make run
```

Local development with auto reload:

```bash
# from go-api/
make dev-install
make dev
```

Build:

```bash
make build
```

## Documentation

- [Docs Index](docs/index.md)
- [Quickstart](docs/quickstart.md)
- [Configuration](docs/configuration.md)
- [Development](docs/development.md)
- [CI/CD Workflow](docs/ci-cd.md)
- [Convention (DO/DON'T)](docs/convention.md)
- [gRPC Testing](docs/grpc-testing.md)
- [Testing](docs/testing.md)
- [API Examples](docs/api-examples.md)
- [Project Structure](docs/project-structure.md)
- [Migrations](docs/MIGRATIONS.md)
- [Kubernetes Deployment Example](docs/k8s-deployment.yaml)

## Proto Contracts

- [User Proto](proto/user/v1/user.proto)
- [Role Proto](proto/role/v1/role.proto)
- [Common Proto](proto/common/v1)
