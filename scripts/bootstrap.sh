#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"
DB_NAME="${DB_NAME:-go_api}"

if [[ "${1:-}" == "--help" ]]; then
  cat <<'EOF'
Usage:
  ./scripts/bootstrap.sh

Optional environment overrides:
  DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME

Example:
  DB_PASSWORD=secret ./scripts/bootstrap.sh
EOF
  exit 0
fi

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "❌ Missing required command: $1"
    exit 1
  fi
}

echo "==> Checking required tools"
require_cmd go
require_cmd protoc
require_cmd make
require_cmd psql
require_cmd sqlc

echo "==> Running code generation and dependency sync"
make proto-tools
make proto
make sqlc
make tidy

echo "==> Applying migration: db/migrations/001_create_users.sql"
PGPASSWORD="$DB_PASSWORD" psql \
  -h "$DB_HOST" \
  -p "$DB_PORT" \
  -U "$DB_USER" \
  -d "$DB_NAME" \
  -f db/migrations/001_create_users.sql

echo "==> Starting service"
make run
