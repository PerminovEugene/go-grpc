# Analytics gRPC Server

A gRPC server that provides analytics services for ticket rating data.

Database file is located in root of `./backend` folder.

## Grpc Service stack

### Prerequisites

- Go 1.21+
- Protocol Buffers compiler (protoc)
- SQLite3

### Local setup

1. Create `.env` file and fill it up, using `.env.example` as template

2. docker compose up -d

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
