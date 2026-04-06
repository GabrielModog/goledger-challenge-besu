# App

This is where your solution goes.

Build your Go application here. Refer to the [challenge instructions](../README.md) for the full requirements, deliverables, and the provided Go code snippets for interacting with the Besu network.

# GoLedger Challenge - API

Blockchain-Database synchronization API using go-fiber and PostgreSQL.

## Architecture Decisions

### 1. go-fiber
- **Great Performance**
- **Concise API**: Less boilerplate for routes and middleware
- **Express-compatible middleware**

### 2. pgx
- **Native connection pool**: Better connection management
- **Context support**: Native cancellation and timeouts
- **Prepared statements**: Optimized performance for repeated queries

### 3. Repository Pattern
- **Separation of concerns**: Handlers don't know database details
- **Testability**: Makes mocking easy for unit tests
- **Maintainability**: Database changes isolated from business logic

### 4. Structured Logging (slog)
- **Go 1.25+ standard**: Native module, no external dependencies
- **JSON output**: Ideal for container logs (Docker/K8s)
- **Log levels**: INFO, WARN, ERROR with structured attributes

### 5. Graceful Shutdown
- **SIGINT/SIGTERM signals**: Clean shutdown without orphaned requests
- **Configurable timeout**: Prevents indefinite hangs
- **Defer resources**: BD/Blockchain connections closed properly

### 6. Multi-stage Dockerfile
- **Small final image**
- **No source code in production**: Only compiled binary
- **Optimized build cache**: Dependencies downloaded before code

## Tech Stack

| Component | Technology |
|-----------|------------|
| Web Framework | go-fiber |
| Blockchain | go-ethereum |
| Database | pgx (Postgres) |
| Config | godotenv |
| Testing | testify |

## Prerequisites

- Go 1.25+
- Docker
- Hyperledger Besu network (QBFT) running
- PostgreSQL (via Docker)

## How to Run

### 1. Initial Setup

```bash
# Enter directory
cd app

# Copy environment variables
cp .env.example .env

# Edit .env with your Besu network values
nano .env
```

### 2. Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `BESU_NODE_URL` | Besu RPC endpoint | `http://127.0.0.1:8545` |
| `CONTRACT_ADDRESS` | SimpleStorage address | `0x...` |
| `SIGNER_PRIVATE_KEY` | Private key for signing TXs | `0x...` |
| `CONTRACT_ABI_PATH` | Path to ABI JSON | `./abi/SimpleStorage.abi.json` |
| `CONNECTION_STRING` | PostgreSQL DSN | `postgres://...` |
| `SERVER_PORT` | Server port (optional) | `8080` |

### 3. Start Infrastructure (Docker)

```bash
# Start PostgreSQL only
docker-compose up -d postgres

# Or start everything (app + postgres)
docker-compose up -d
```

### 4. Run Locally (without Docker)

```bash
# Download dependencies
go mod tidy

# Run
go run cmd/server/main.go
```

### 5. Build Binary

```bash
go build -o goledger-server ./cmd/server
./goledger-server
```

## API Endpoints

| Method | Endpoint | Description | Body |
|--------|----------|-------------|------|
| GET | `/health` | Health check | - |
| GET | `/api/v1/get` | Read value from blockchain | - |
| POST | `/api/v1/set` | Set value on blockchain | `{"value": "42"}` |
| POST | `/api/v1/sync` | Sync BC → PostgreSQL | - |
| GET | `/api/v1/check` | Check consistency | - |

### curl Examples

```bash
# Set value
curl -X POST http://localhost:8080/api/v1/set \
  -H "Content-Type: application/json" \
  -d '{"value": "42"}'

# Get value
curl http://localhost:8080/api/v1/get

# Sync to database
curl -X POST http://localhost:8080/api/v1/sync

# Check consistency
curl http://localhost:8080/api/v1/check
```

## Development

### Run Tests

```bash
go test ./... -v
```

### Logs

```bash
# Docker
docker-compose logs -f app

# Local
# Logs are output in JSON to stdout
```

## Container Healthcheck

The Dockerfile includes a healthcheck that verifies `/health`:

```yaml
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1
```
