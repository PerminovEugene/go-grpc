.PHONY: proto build run clean test test-watch

# Generate protobuf code
SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c

PROTO_DIR := proto

.PHONY: proto
proto:
	cd $(PROTO_DIR) && \
	protoc -I . \
	  --go_out=paths=source_relative:. \
	  --go-grpc_out=paths=source_relative:. \
	  $$(find . -name '*.proto' -print)



# Build the application
build:
	go build -o bin/server main.go

# Build the client
build-client:
	go build -o bin/client client

# Run the server
run:
	go run main.go

# Run the client (make sure server is running first)
run-client:
	go run client

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
deps:
	go mod tidy
	go mod download

# Install protoc dependencies (run this first)
install-protoc:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Run tests
test:
	go test ./... -v

# Run tests in watch mode
test-watch:
	go run cmd/test-watcher/main.go
