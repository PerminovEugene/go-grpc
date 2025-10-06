# Analytics gRPC Server

A gRPC server that provides analytics services for ticket rating data.

## Prerequisites

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

Preferably to put ./database.db file to ./backend folder

Available environment variables:
- `DB_PATH` - Path to SQLite database file (default: `database.db`)
- `IMAGE_TAG` - Docker image tag for versioning (default: `latest`)
- `GRPC_PORT` - gRPC server port (default: `50051`)

### 2. Docker Deployment

**Build and run with default settings:**
```bash
docker compose up -d
```

**Build with a custom image tag:**
```bash
IMAGE_TAG=v1.0.0 docker compose up -d
```

**Or set it in your .env file:**
```bash
echo "IMAGE_TAG=v1.0.0" >> .env
docker compose up -d
```

### 3. Local Development (without Docker)

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

1. make test

## Production

1. TODO

## Database

The server uses SQLite database, the file with db data is in ./backend folder

## API

### GetAggregatedCategoryScores

Returns daily aggregates for periods ≤ 1 month, weekly for longer periods.

### GetScoresByTicket

Returns scores grouped by ticket within a period.

### GetOverallQualityScore

Returns overall quality score for a period.

### GetPeriodOverPeriodChange

Compares two periods and shows percentage change.
