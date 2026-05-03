# Testing

## Unit tests

```bash
# tests under test/
go test ./test -v

# all tests in module
go test ./... -v
```

## API script test (gRPC smoke test)

Run full API smoke tests and generate request/response docs:

```bash
# from go-api root
bash scripts/test_all_apis.sh
```

Generated output:

- `docs/api_test_results.md`

### Prerequisites

- gRPC server is running at `localhost:50051`
- Database has seed data used by script (user/product/module IDs in `scripts/test_all_apis.sh`)
- `grpcurl`, `python3`, and `psql` are available in your shell

### Typical workflow

```bash
# 1) start server
make build && ./bin/grpc_server

# 2) in another terminal, run smoke test script
bash scripts/test_all_apis.sh
```

### Troubleshooting

If server fails with `bind: address already in use`:

```bash
lsof -i :50051
kill <PID>
```

## Coverage

Generate coverage file:

```bash
go test ./test -coverprofile=coverage.out
```

Show summary:

```bash
go tool cover -func=coverage.out
```

Generate HTML report:

```bash
go tool cover -html=coverage.out -o coverage.html
```
