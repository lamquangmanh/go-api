APP_NAME=go-api
CONFIG ?= configs/config.yml
GO_BIN := $(shell if [ -n "$$(go env GOBIN)" ]; then echo "$$(go env GOBIN)"; else echo "$$(go env GOPATH)/bin"; fi)
AIR_BIN := $(GO_BIN)/air

.PHONY: run dev dev-install build docker sqlc migrate migrate-tool-install migrate-up migrate-down migrate-status proto proto-tools tidy fmt fmt-check

run:
	GOAPI_CONFIG=$(CONFIG) go run ./cmd/grpc

dev:
	@if command -v air >/dev/null 2>&1; then \
		air -c .air.toml; \
	elif [ -x "$(AIR_BIN)" ]; then \
		"$(AIR_BIN)" -c .air.toml; \
	else \
		echo "Error: air not found"; \
		echo "Run: make dev-install"; \
		echo "Or add $(GO_BIN) to PATH"; \
		exit 1; \
	fi

dev-install:
	go install github.com/air-verse/air@latest
	@echo "air installed to: $(AIR_BIN)"
	@echo "If needed, add to PATH: export PATH=\"$(GO_BIN):$$PATH\""

build:
	go build -o bin/grpc ./cmd/grpc

# Format all Go code files using gofmt with -w flag (write formatted code back to files)
# This ensures consistent code formatting across the entire project
# Usage: make fmt
fmt:
	@echo "Formatting all Go code files..."
	@gofmt -l -w ./cmd ./internal ./pkg
	@echo "✓ All Go files formatted successfully"

# Check Go code formatting without making changes (useful for CI/CD validation)
# This verifies that all Go files follow the standard gofmt style
# Fails if any files need formatting - typically used in GitHub Actions workflows
# Usage: make fmt-check
fmt-check:
	@echo "Checking Go code formatting..."
	@if [ -n "$$(gofmt -l ./cmd ./internal ./pkg)" ]; then \
		echo "❌ The following files need formatting:"; \
		gofmt -l ./cmd ./internal ./pkg; \
		echo "Run 'make fmt' to auto-format or use 'gofmt -w <file>' manually"; \
		exit 1; \
	fi
	@echo "✓ All Go files are properly formatted"

docker:
	docker build -t $(APP_NAME):latest .

sqlc:
	sqlc generate

# Install golang-migrate CLI tool
# Usage: make migrate-tool-install
migrate-tool-install:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "golang-migrate installed successfully"
	@echo "Usage: migrate -path db/migrations -database 'postgres://user:pass@localhost:5432/dbname?sslmode=disable' up"

# Apply all pending migrations
# Usage: make migrate-up DB_URL='postgres://user:pass@localhost:5432/dbname?sslmode=disable'
migrate-up:
	@if [ -z "$(DB_URL)" ]; then \
		echo "Error: DB_URL not set"; \
		echo "Usage: make migrate-up DB_URL='postgres://user:pass@localhost:5432/dbname?sslmode=disable'"; \
		exit 1; \
	fi
	migrate -path db/migrations -database "$(DB_URL)" up

# Rollback one migration step
# Usage: make migrate-down DB_URL='postgres://user:pass@localhost:5432/dbname?sslmode=disable'
migrate-down:
	@if [ -z "$(DB_URL)" ]; then \
		echo "Error: DB_URL not set"; \
		echo "Usage: make migrate-down DB_URL='postgres://user:pass@localhost:5432/dbname?sslmode=disable'"; \
		exit 1; \
	fi
	migrate -path db/migrations -database "$(DB_URL)" down 1

# Show current migration version
# Usage: make migrate-status DB_URL='postgres://user:pass@localhost:5432/dbname?sslmode=disable'
migrate-status:
	@if [ -z "$(DB_URL)" ]; then \
		echo "Error: DB_URL not set"; \
		echo "Usage: make migrate-status DB_URL='postgres://user:pass@localhost:5432/dbname?sslmode=disable'"; \
		exit 1; \
	fi
	@echo "Current migration version:"
	@migrate -path db/migrations -database "$(DB_URL)" version

# Legacy target for backward compatibility - applies all migrations
# Deprecated: Use make migrate-up instead with DB_URL parameter
migrate:
	@echo "⚠️  'make migrate' is deprecated. Use 'make migrate-up DB_URL=<db_url>' instead."
	@if [ -z "$(DB_URL)" ]; then \
		echo "Example: make migrate-up DB_URL='postgres://user:pass@localhost:5432/dbname?sslmode=disable'"; \
		exit 1; \
	fi
	@make migrate-up DB_URL=$(DB_URL)

proto-tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

proto:
	PATH="$(PATH):$(shell go env GOPATH)/bin" protoc -I .. --go_out=. --go_opt=module=go-api --go-grpc_out=. --go-grpc_opt=module=go-api $(shell find proto -name "*.proto" | sed 's#^#go-api/#')

tidy:
	go mod tidy

