.PHONY: proto build build-debug debug run clean test test-watch install-delve fmt

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

# Build with debug symbols for delve
build-debug:
	go build -gcflags="all=-N -l" -o bin/server-debug main.go

# Build the category scores client
build-client:
	go build -o bin/category_scores_client ./client/category_scores

# Build the ticket scores client
build-ticket-client:
	go build -o bin/ticket_scores_client ./client/ticket_scores

# Build the overall quality score client
build-overall-quality-client:
	go build -o bin/overall_quality_score_client ./client/overall_quality_score

# Build the period over period client
build-period-over-period-client:
	go build -o bin/period_over_period_client ./client/period_over_period

# Run the server
run:
	go run main.go

# Run the server with delve debugger (headless mode, connect with VSCode)
# Usage: make debug [PORT=2345]
debug:
	@echo "Starting delve debugger on port $(or $(PORT),2345)..."
	@echo "Attach VSCode debugger or use: dlv connect localhost:$(or $(PORT),2345)"
	dlv debug main.go --headless --listen=:$(or $(PORT),2345) --api-version=2 --accept-multiclient

# Run the category scores client (make sure server is running first)
# Usage: make run-client START=2025-01-01 END=2025-01-31
# Or: make run-client (uses default last 30 days)
run-client:
	@if [ -z "$(START)" ] && [ -z "$(END)" ]; then \
		go run ./client/category_scores; \
	elif [ -z "$(START)" ] || [ -z "$(END)" ]; then \
		echo "Error: Both START and END must be provided together"; \
		echo "Usage: make run-client START=2025-01-01 END=2025-01-31"; \
		exit 1; \
	else \
		go run ./client/category_scores -start $(START) -end $(END); \
	fi

# Run the ticket scores client (make sure server is running first)
# Usage: make run-ticket-client START=2025-01-01 END=2025-01-31
# Or: make run-ticket-client (uses default last 30 days)
run-ticket-client:
	@if [ -z "$(START)" ] && [ -z "$(END)" ]; then \
		go run ./client/ticket_scores; \
	elif [ -z "$(START)" ] || [ -z "$(END)" ]; then \
		echo "Error: Both START and END must be provided together"; \
		echo "Usage: make run-ticket-client START=2025-01-01 END=2025-01-31"; \
		exit 1; \
	else \
		go run ./client/ticket_scores -start $(START) -end $(END); \
	fi

# Run the overall quality score client (make sure server is running first)
# Usage: make run-overall-quality-client START=2025-01-01 END=2025-01-31
# Or: make run-overall-quality-client (uses default last 7 days)
run-overall-quality-client:
	@if [ -z "$(START)" ] && [ -z "$(END)" ]; then \
		go run ./client/overall_quality_score; \
	elif [ -z "$(START)" ] || [ -z "$(END)" ]; then \
		echo "Error: Both START and END must be provided together"; \
		echo "Usage: make run-overall-quality-client START=2025-01-01 END=2025-01-31"; \
		exit 1; \
	else \
		go run ./client/overall_quality_score -start $(START) -end $(END); \
	fi

# Run the period over period client (make sure server is running first)
# Usage: make run-period-over-period-client CURR_START=2025-02-01 CURR_END=2025-02-28 PREV_START=2025-01-01 PREV_END=2025-01-31
# Or: make run-period-over-period-client (uses default: current week vs previous week)
run-period-over-period-client:
	@if [ -z "$(CURR_START)" ] && [ -z "$(CURR_END)" ] && [ -z "$(PREV_START)" ] && [ -z "$(PREV_END)" ]; then \
		go run ./client/period_over_period; \
	elif [ -z "$(CURR_START)" ] || [ -z "$(CURR_END)" ] || [ -z "$(PREV_START)" ] || [ -z "$(PREV_END)" ]; then \
		echo "Error: All four dates must be provided together"; \
		echo "Usage: make run-period-over-period-client CURR_START=2025-02-01 CURR_END=2025-02-28 PREV_START=2025-01-01 PREV_END=2025-01-31"; \
		exit 1; \
	else \
		go run ./client/period_over_period -current-start $(CURR_START) -current-end $(CURR_END) -previous-start $(PREV_START) -previous-end $(PREV_END); \
	fi

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
deps:
	go mod tidy
	go mod download

# Format Go code
fmt:
	go fmt ./...
	gofmt -s -w .

# Install protoc dependencies (run this first)
install-protoc:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Install delve debugger
install-delve:
	go install github.com/go-delve/delve/cmd/dlv@latest

# Run tests
test:
	go test ./... -v

# Run tests in watch mode
test-watch:
	go run cmd/test-watcher/main.go
