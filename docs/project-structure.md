# Project Structure

```text
go-api/
├── cmd/grpc/                 # server entrypoint
├── configs/                  # runtime config
├── internal/
│   ├── config/               # config loader
│   ├── db/                   # pgxpool helpers
│   ├── grpc/                 # transport handlers
│   ├── logger/               # structured logger + interceptor
│   ├── repository/           # sqlc-generated code
│   └── service/              # business logic
├── db/
│   ├── migrations/           # golang-migrate up/down SQL
│   └── query/                # sqlc query definitions
├── pkg/
│   ├── api/                  # generated protobuf stubs
│   ├── constants/            # shared error constants
│   └── utils/                # shared utilities
├── proto/                    # protobuf contracts
├── test/                     # tests
└── docs/                     # detailed docs
```
