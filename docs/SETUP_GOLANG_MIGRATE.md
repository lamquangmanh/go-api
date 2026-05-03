# golang-migrate Setup Summary

## ✅ Completed Tasks

### 1. Migration File Format Conversion

- ✅ Converted 3 migration files from single-file format to golang-migrate up/down format:
  - `000001_init_users.up.sql` + `.down.sql`
  - `000002_create_roles.up.sql` + `.down.sql`
  - `000003_add_audit_columns.up.sql` + `.down.sql`

**Location:** `db/migrations/`

**Benefits:**

- ✅ Atomic rollback support (down migrations can undo changes)
- ✅ Version tracking in `schema_migrations` table
- ✅ Prevents duplicate runs automatically
- ✅ Supports dirty state detection and recovery

### 2. Dependency Added

- ✅ Added `github.com/golang-migrate/migrate/v4` to `go.mod`
- ✅ Project builds successfully with new dependency

### 3. Makefile Commands

- ✅ Added `make migrate-up` - Apply all pending migrations
- ✅ Added `make migrate-down` - Rollback one migration
- ✅ Added `make migrate-status` - Check current migration version
- ✅ Added `make migrate-tool-install` - Install CLI tool locally

**Usage Examples:**

```bash
# Install the tool
make migrate-tool-install

# Apply migrations
make migrate-up DB_URL='postgresql://user:pass@localhost:5432/dbname?sslmode=disable'

# Check status
make migrate-status DB_URL='postgresql://user:pass@localhost:5432/dbname?sslmode=disable'

# Rollback if needed
make migrate-down DB_URL='postgresql://user:pass@localhost:5432/dbname?sslmode=disable'
```

### 4. Docker & Kubernetes Integration

- ✅ Updated `Dockerfile` to:
  - Build `migrate` binary with PostgreSQL support
  - Include migration files in runtime image
- ✅ Created `docker-compose.yml` with:
  - PostgreSQL service
  - Optional migration service
  - Go API application
  - PgAdmin for debugging (optional)
  - Testing/development profiles

- ✅ Created `docs/k8s-deployment.yaml` (Kubernetes manifest) with:
  - Init container pattern for database migrations
  - Proper environment variables and secrets management
  - Health checks and resource limits
  - Pod anti-affinity for high availability
  - PodDisruptionBudget for reliable deployments

### 5. Documentation

- ✅ Created `docs/MIGRATIONS.md` (Comprehensive guide):
  - Installation and setup
  - Usage examples for all commands
  - Creating new migrations
  - Best practices for production
  - Zero-downtime deployment strategies
  - Troubleshooting and rollback procedures
  - Kubernetes integration details

- ✅ Updated `README.md` with:
  - New Migrations section
  - Quick start commands
  - Links to comprehensive migration docs
  - Updated schema documentation

---

## 📁 File Structure

**New files created:**

```
go-api/
├── db/migrations/
│   ├── 000001_init_users.up.sql
│   ├── 000001_init_users.down.sql
│   ├── 000002_create_roles.up.sql
│   ├── 000002_create_roles.down.sql
│   ├── 000003_add_audit_columns.up.sql
│   └── 000003_add_audit_columns.down.sql
├── docs/
│   ├── MIGRATIONS.md (comprehensive guide)
│   └── k8s-deployment.yaml (kubernetes manifest)
├── docker-compose.yml (development setup)
└── go.mod (updated with golang-migrate)
```

**Updated files:**

- `Makefile` - New migration commands
- `Dockerfile` - Build and include migrate binary
- `README.md` - Migration documentation and examples

---

## 🚀 Quick Start

### Local Development

**Option 1: Using Docker Compose**

```bash
# Start database and apply migrations
docker-compose --profile with-migrations up

# Or use docker-compose without profile for manual control
docker-compose up postgres
```

**Option 2: Manual using Makefile**

```bash
# Install golang-migrate CLI
make migrate-tool-install

# Start PostgreSQL (separately in another terminal or use docker)
docker run -d -p 5432:5432 \
  -e POSTGRES_USER=go_api_user \
  -e POSTGRES_PASSWORD=go_api_password \
  -e POSTGRES_DB=go_api \
  postgres:16

# Apply migrations
make migrate-up DB_URL='postgresql://go_api_user:go_api_password@localhost:5432/go_api?sslmode=disable'

# Build and run the service
make run CONFIG=configs/config.yml
```

### Kubernetes Deployment

```bash
# 1. Build Docker image with migrations included
docker build -t go-api:v1.0 .
docker push your-registry/go-api:v1.0

# 2. Apply Kubernetes manifests
kubectl apply -f docs/k8s-deployment.yaml

# 3. Verify migrations ran
kubectl logs -f deployment/go-api -c db-migrate

# 4. Check application status
kubectl get pods -l app=go-api
kubectl logs -f deployment/go-api -c grpc-server
```

---

## 📚 Key Concepts

### Migration Lifecycle

```
1. Developer writes *.up.sql and *.down.sql
   ↓
2. Commit to git
   ↓
3. Docker build includes migration files
   ↓
4. Deployment starts (dev/prod):
   - Dev: `make migrate-up DB_URL=...`
   - K8s: Init container runs `migrate up` before app container
   ↓
5. Migrations tracked in `schema_migrations` table
   ↓
6. If rollback needed: `make migrate-down` or K8s redeploy
```

### Kubernetes Pattern

```
Pod Lifecycle:
1. Init Container: Run migrations (db-migrate)
2. Readiness Probe: Checks db-migrate status
3. App Container: Starts only after migrations complete
4. Liveness Probe: Ensures app stays healthy
```

**Advantages:**

- ✅ No manual migration step before deployment
- ✅ Atomic: All or nothing (failures block app startup)
- ✅ Rollback-friendly: Can redeploy previous version
- ✅ Scalable: Works with multipl pod replicas

---

## ⚠️ Important Notes

### Old Migration Files

Original files are kept for reference but no longer used:

- `001_create_users.sql` → Use `000001_init_users.up.sql` instead
- `002_create_roles_and_user_role_mappings.sql` → Use `000002_create_roles.up.sql`
- `003_add_audit_columns.sql` → Use `000003_add_audit_columns.up.sql`

To clean up (optional):

```bash
cd db/migrations
rm 001_*.sql 002_*.sql 003_*.sql  # Keep only 000001, 000002, 000003
```

### Schema Migrations Table

golang-migrate automatically creates and manages:

```sql
CREATE TABLE schema_migrations (
    version BIGINT PRIMARY KEY,
    dirty   BOOLEAN NOT NULL
);
```

**Never modify this table manually** unless recovering from a failed migration.

---

## 🔄 Workflow Examples

### Create a New Migration

```bash
# Tool will ask for migration name
cd /Users/jun/Documents/Projects/lamquangmanh-github/go-api
migrate create -ext sql -dir db/migrations -seq add_user_status_column

# Creates: db/migrations/000004_add_user_status_column.up.sql
#          db/migrations/000004_add_user_status_column.down.sql

# Edit both files with your SQL
# Then apply:
make migrate-up DB_URL='postgresql://...'
```

### Deploy New Version

```bash
# 1. Add migration files to db/migrations/
# 2. Commit and push
# 3. Build Docker image (includes migrations)
docker build -t go-api:v1.1 .

# 4. Update K8s manifest or trigger deployment
kubectl set image deployment/go-api grpc-server=go-api:v1.1

# 5. Kubernetes handles the rest:
# - Init container runs new migrations automatically
# - App container starts after migrations succeed
# - Pods scale according to PDB and health checks
```

### Rollback on Issues

```bash
# Quick rollback using previous K8s deployment
kubectl rollout undo deployment/go-api

# Manual rollback (if needed):
make migrate-down DB_URL='postgresql://...'

# Redeploy previous version
kubectl set image deployment/go-api grpc-server=go-api:v1.0
```

---

## 📖 Further Reading

- [golang-migrate Documentation](https://github.com/golang-migrate/migrate)
- [docs/MIGRATIONS.md](MIGRATIONS.md) - Comprehensive migration guide
- [docs/k8s-deployment.yaml](k8s-deployment.yaml) - Full K8s example
- [docker-compose.yml](../docker-compose.yml) - Local development setup

---

## ✨ Next Steps (Optional)

1. **Test migrations locally** with `docker-compose up`
2. **Deploy to staging** to verify K8s workflow
3. **Set up CI/CD** to validate migrations in PR/MR pipelines
4. **Document** any custom migration scripts or procedures
5. **Monitor** migration duration in production deployments
