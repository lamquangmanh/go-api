# Database Migration Guide

## Overview

This project uses [golang-migrate](https://github.com/golang-migrate/migrate) for database schema versioning and management. Migrations are organized in the `db/migrations/` directory with up/down SQL files for atomic rollback support.

## File Structure

```
db/migrations/
├── 000001_init_users.up.sql
├── 000001_init_users.down.sql
├── 000002_create_roles.up.sql
├── 000002_create_roles.down.sql
├── 000003_add_audit_columns.up.sql
└── 000003_add_audit_columns.down.sql
```

### Naming Convention

- **Pattern**: `{VERSION}_{DESCRIPTION}.{DIRECTION}.sql`
- **VERSION**: Zero-padded sequence (000001, 000002, etc.)
- **DESCRIPTION**: Descriptive name in snake_case
- **DIRECTION**: `up` (apply) or `down` (rollback)

**Examples:**

- `000001_init_users.up.sql` - Create users table
- `000001_init_users.down.sql` - Drop users table
- `000002_create_roles.up.sql` - Create roles table
- `000002_create_roles.down.sql` - Drop roles table

## Installation

### Option 1: Install CLI Tool Locally

```bash
# Install golang-migrate CLI
make migrate-tool-install

# Verify installation
migrate -version
```

### Option 2: Using Docker

```bash
docker run --rm -v $(pwd)/db/migrations:/migrations migrate/migrate \
  -path=/migrations \
  -database "postgresql://user:password@localhost:5432/dbname?sslmode=disable" \
  up
```

## Usage

### Environment Setup

Create a `.env` file or export the database URL:

```bash
# PostgreSQL connection string
export DB_URL="postgresql://go_api_user:go_api_password@localhost:5432/go_api?sslmode=disable"
```

### Apply All Pending Migrations (Up)

```bash
# Using Makefile
make migrate-up DB_URL='postgresql://user:pass@localhost:5432/dbname?sslmode=disable'

# Using CLI directly
migrate -path db/migrations -database "postgresql://user:pass@localhost:5432/dbname?sslmode=disable" up

# Using Docker
docker run --rm -v $(pwd)/db/migrations:/migrations \
  migrate/migrate -path=/migrations \
  -database "postgresql://user:pass@host:5432/dbname?sslmode=disable" \
  up
```

### Rollback One Migration (Down)

```bash
# Using Makefile
make migrate-down DB_URL='postgresql://user:pass@localhost:5432/dbname?sslmode=disable'

# Using CLI directly
migrate -path db/migrations -database "postgresql://user:pass@localhost:5432/dbname?sslmode=disable" down 1
```

### Check Current Migration Version

```bash
# Using Makefile
make migrate-status DB_URL='postgresql://user:pass@localhost:5432/dbname?sslmode=disable'

# Using CLI directly
migrate -path db/migrations -database "postgresql://user:pass@localhost:5432/dbname?sslmode=disable" version
```

### Force-Set Migration Version (Emergency Only)

```bash
# Force set to version N (use with caution!)
migrate -path db/migrations -database "postgresql://..." force 3

# Check dirty state
migrate -path db/migrations -database "postgresql://..." version
```

## Creating New Migrations

### Step 1: Create Migration Files

```bash
# Using golang-migrate CLI
migrate create -ext sql -dir db/migrations -seq add_user_status_column

# This creates:
# db/migrations/000004_add_user_status_column.up.sql
# db/migrations/000004_add_user_status_column.down.sql
```

### Step 2: Write SQL

**000004_add_user_status_column.up.sql:**

```sql
-- Add status column to users table
ALTER TABLE users
    ADD COLUMN status VARCHAR(50) DEFAULT 'active',
    ADD COLUMN status_updated_at TIMESTAMPTZ DEFAULT NOW();

CREATE INDEX idx_users_status ON users(status);
```

**000004_add_user_status_column.down.sql:**

```sql
-- Rollback: Remove status column from users table
DROP INDEX IF EXISTS idx_users_status;

ALTER TABLE users
    DROP COLUMN IF EXISTS status,
    DROP COLUMN IF EXISTS status_updated_at;
```

### Step 3: Test Migration

```bash
# Apply migration
make migrate-up DB_URL='postgresql://...'

# Test your application

# Rollback if issues found
make migrate-down DB_URL='postgresql://...'
```

## Kubernetes Deployment

### Init Container Pattern

Use init containers to apply migrations before app startup:

```yaml
initContainers:
  - name: db-migrate
    image: go-api:latest
    command:
      - /app/migrate
      - -path=/app/db/migrations
      - -database=$(DATABASE_URL)
      - up
    env:
      - name: DATABASE_URL
        valueFrom:
          secretKeyRef:
            name: db-secret
            key: url
```

See [docs/k8s-deployment.yaml](k8s-deployment.yaml) for complete example.

### Deployment Process

```bash
# 1. Build Docker image with migrations included
docker build -t go-api:v1.2.3 .

# 2. Push to registry
docker push registry.example.com/go-api:v1.2.3

# 3. Deploy to Kubernetes
kubectl apply -f docs/k8s-deployment.yaml

# 4. Verify migrations ran
kubectl logs -f deployment/go-api -c db-migrate
```

## IMPORTANT CONSIDERATIONS

### ⚠️ Never Break Compatibility

- **Down migrations must always work** - You'll need to rollback for hotfixes
- **Test down migrations** - Not just up!
- **Avoid destructive changes without down script** - Example: dropping columns without reversibility

### ⚠️ Zero-Downtime Deployments

For large tables, use **safe migration practices:**

```sql
-- ✅ GOOD: Add column with default, migrate data, then add constraint
ALTER TABLE users ADD COLUMN email_verified BOOLEAN DEFAULT false;
UPDATE users SET email_verified = true WHERE verified_at IS NOT NULL;
ALTER TABLE users ALTER COLUMN email_verified SET NOT NULL;

-- ❌ BAD: Immediate breaking changes
ALTER TABLE users DROP COLUMN status CASCADE;
```

### ⚠️ High Volume Schema Changes

For tables with millions of rows:

1. Use `CONCURRENTLY` for index creation (separate migration)
2. Add columns with defaults (non-blocking)
3. Use batched updates for data migration
4. Drop old columns in follow-up migration (after verification)

Example:

```sql
-- Safe for large tables
CREATE INDEX CONCURRENTLY idx_users_status ON users(status);
-- Wrap in try-catch in app code to handle concurrent creation
```

## Monitoring Migration Status

### Check Applied Migrations

```bash
psql $DB_URL -c "SELECT * FROM schema_migrations ORDER BY version DESC;"
```

**Output:**

```
 version |                dirty |
---------+----------------------+
 3       | f                    |
 2       | f                    |
 1       | f                    |
```

### Handle Dirty State

If a migration fails, the `dirty` flag is set. **Fix it manually:**

```bash
# 1. Check what went wrong
psql $DB_URL -c "SELECT * FROM schema_migrations WHERE dirty = true;"

# 2. Manually fix database state if needed

# 3. Reset dirty flag
make migrate-force DB_VERSION=3  # Force to version 3 (or last good version)

# 4. Retry migration
make migrate-up DB_URL='...'
```

## Best Practices

### 1. Version Control Migrations

```bash
# Always commit migration files to git
git add db/migrations/
git commit -m "chore: add migration 000004_add_user_status_column"
```

### 2. Review Migration PRs

- Verify SQL syntax and correctness
- Check that down migration also works
- Test against production data subset locally
- Review for data loss risks

### 3. Local Development

```bash
# Start fresh database
docker-compose up -d postgres

# Apply all migrations
make migrate-up DB_URL='...'

# Run tests
go test ./test/...
```

### 4. Production Deployment

```bash
# Before deploying new app version:
# 1. Run migrations on staging
make migrate-up DB_URL=$STAGING_DB_URL

# 2. Test application on staging
# 3. Get approval
# 4. Deploy to production (init container handles migration)
```

### 5. Rollback Procedure

```bash
# If deployment fails:
# 1. Check pod logs
kubectl logs -f deployment/go-api -c db-migrate

# 2. Rollback kubernetes deployment
kubectl rollout undo deployment/go-api

# 3. If data corruption, manually fix via DB access
# 4. Re-apply migrations after fix:
make migrate-down DB_URL=$PROD_DB_URL  # Rollback last migration
make migrate-up DB_URL=$PROD_DB_URL    # Try again
```

## Troubleshooting

### Error: "Dirty database state"

```
error: Dirty database
```

**Solution:**

```bash
# 1. Check what migration failed
psql $DB_URL -c "SELECT * FROM schema_migrations WHERE dirty = true;"

# 2. Manually fix the database state

# 3. Force reset
migrate -path db/migrations -database "$DB_URL" force 2  # Force to version 2

# 4. Verify clean state
migrate -path db/migrations -database "$DB_URL" version
```

### Error: "Version doesn't exist"

```
error: Unable to determine dirty state
```

**Causes:**

- Connection timeout
- schema_migrations table corrupted
- Connection string error

**Solution:**

```bash
# Verify connection
psql $DB_URL -c "SELECT 1"

# Check schema_migrations table
psql $DB_URL -c "SELECT * FROM schema_migrations;"

# Recreate if corrupted
psql $DB_URL -c "DROP TABLE IF EXISTS schema_migrations;"
make migrate-up DB_URL=$DB_URL  # Will recreate table
```

### Error: "No change" during migration

Possible causes:

- Idempotent migration (uses `IF NOT EXISTS`) - this is OK
- Previous partial migration - check dirty state
- Migration already applied - check version

## References

- [golang-migrate Documentation](https://github.com/golang-migrate/migrate)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Kubernetes Init Containers](https://kubernetes.io/docs/concepts/workloads/pods/init-containers/)
