FROM golang:1.23.2-alpine AS builder

WORKDIR /app

# Install build dependencies for SQLite
RUN apk add --no-cache gcc musl-dev sqlite-dev

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN go build -o server main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates sqlite

WORKDIR /app

# Accept DB_PATH as build argument with default value
ARG DB_PATH=./database.db

COPY --from=builder /app/server .
COPY --from=builder /app/${DB_PATH} ./database.db

CMD ["./server"]
