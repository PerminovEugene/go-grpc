# Analytics gRPC Server

A gRPC server that provides analytics services for ticket rating data.

All service code is in `./backend` folder.

## Stack

- Go 1.21+
- Protocol Buffers compiler (protoc)
- SQLite3
- Docker and Docker Compose (for containerized deployment)

## Local Setup

### 1. Environment Configuration

Create a `.env` file in the project root using `.env.example` as a template:

```bash
cp .env.example .env
```

Put ./database.db file to ./backend folder

Available environment variables:
- `DB_PATH` - Path to SQLite database file (default: `database.db`)
- `GRPC_PORT` - gRPC server port (default: `50051`)


### 2. Local Development (without Docker)

Run the server directly:
```bash
cd backend
make run
```

Build the server:
```bash
cd backend
make build
```

## Testing

1. `make test`

## CI/CD

Docker images are automatically built and published to GitHub Container Registry on:
- Push to `main` branch → `latest` tag
- Version tags (e.g., `v1.0.0`) → semantic version tags

To create a new release:
```bash
git tag v1.0.0
git push origin v1.0.0
```

See [.github/workflows/README.md](.github/workflows/README.md) for more details.

## Production

Docker images are available at `ghcr.io/PerminovEugene/go-grpc` and can be deployed to any container orchestration platform (Kubernetes, Docker Swarm, etc.).

## API

### GetAggregatedCategoryScores

Returns daily aggregates for periods ≤ 1 month, weekly for longer periods.

### GetScoresByTicket

Returns scores grouped by ticket within a period.

### GetOverallQualityScore

Returns overall quality score for a period.

### GetPeriodOverPeriodChange

Compares two periods and shows percentage change.
