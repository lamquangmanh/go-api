# =============================================================================
# STAGE 1: Builder
# Use the golang:1.25 image as the build environment.
# This stage compiles the Go source code into a binary artifact.
# It is named "builder" so later stages can reference and copy the built binary.
# =============================================================================
FROM golang:1.25 AS builder

# Set the working directory inside the container to /app.
# Subsequent COPY and RUN commands execute from this directory.
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker layer caching.
# Downloading dependencies only runs again if these files change,
# which speeds up rebuilds when only application code is modified.
COPY go.mod go.sum ./

# Download all modules required by the project into the module cache.
RUN go mod download

# Copy the rest of the source code into the container after dependencies are cached.
COPY . .

# Build the gRPC server binary:
#   CGO_ENABLED=0  — disable CGO to produce a static binary without external C libs.
#   GOOS=linux     — target operating system is Linux.
#   GOARCH=amd64   — target architecture x86-64.
#   -o grpc        — output binary named `grpc`.
#   ./cmd/grpc     — package containing main for the gRPC server.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o grpc ./cmd/grpc

# Build the database migration tool with PostgreSQL support:
# This is used by the init container to apply migrations before app startup.
# Note: go install is used because the migrate module is not in go.mod.
# It automatically downloads and compiles the tool in the current directory.
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest && \
    cp /go/bin/migrate ./migrate

# =============================================================================
# STAGE 2: Runtime (gRPC server)
# Use a distroless base image for a minimal runtime without shell or package managers.
# This reduces image size and surface area for potential vulnerabilities.
# =============================================================================
FROM gcr.io/distroless/base-debian11

# Set working directory for the runtime image.
WORKDIR /app

# Copy only the compiled binary from the builder stage into the runtime image.
# The final image will not include the Go toolchain or the source code,
# which keeps it small and production-friendly.
COPY --from=builder /app/grpc .
COPY --from=builder /app/migrate .
COPY db/migrations ./db/migrations

# Declare the ports the container will listen on:
#   50051 — default gRPC port.
#   8080  — optional HTTP port (e.g., health checks or gRPC-Gateway).
EXPOSE 50051
EXPOSE 8080

# Default command to run when the container starts: execute the gRPC server binary.
CMD ["./grpc"]
