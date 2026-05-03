# Development

## Local dev with auto reload (Air)

Run once on a new machine:

```bash
make dev-install
```

Start local development server (watch + rebuild + restart on file changes):

```bash
make dev
```

Notes:

- `make dev` uses `.air.toml` in the project root.
- If `air` is not available in your shell `PATH`, `make dev` automatically falls back to the Go bin directory (`GOBIN` or `GOPATH/bin`).
- To run `air` directly from shell, add Go bin to your `PATH`:

```bash
export PATH="$(go env GOPATH)/bin:$PATH"
```

## Useful commands

```bash
make proto        # generate protobuf Go code
make sqlc         # generate repository code from SQL
make tidy         # sync go.mod/go.sum
make run          # run service without auto reload
make dev-install  # install air
make dev          # run with auto reload via air
make build        # build binary to bin/grpc
make docker       # build Docker image
```

## Validation behavior

- `CreateUser` validates username/email format and uniqueness
- `UpdateUser` returns gRPC status error on validation failure
- Shared validators are implemented in `pkg/utils`
- Centralized error definitions are in `pkg/constants/errors.go`

## Logging

- Logger uses `log/slog`
- Default format: JSON
- Request logs include method, code, duration and peer
- `x-request-id` metadata is included when available
