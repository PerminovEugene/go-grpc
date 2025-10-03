.PHONY: proto build run clean

# Generate protobuf code
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/analytics.proto

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
