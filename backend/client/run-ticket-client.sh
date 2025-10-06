#!/bin/bash

# Quick helper script to run the ticket scores client
# Usage examples:
#   ./run-ticket-client.sh 2025-01-01 2025-01-31
#   ./run-ticket-client.sh 2025-01-01 2025-01-31 localhost:50051

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$(dirname "$SCRIPT_DIR")"
CLIENT_BIN="$BACKEND_DIR/bin/ticket_scores_client"

# Check if binary exists, if not build it
if [ ! -f "$CLIENT_BIN" ]; then
    echo "Client binary not found. Building..."
    cd "$BACKEND_DIR" && go build -o bin/ticket_scores_client ./client/ticket_scores
    if [ $? -ne 0 ]; then
        echo "Failed to build client"
        exit 1
    fi
    echo "Build successful!"
    echo ""
fi

# Parse arguments
if [ $# -eq 0 ]; then
    # No arguments - use defaults
    echo "Using default dates (last 30 days)"
    "$CLIENT_BIN"
elif [ $# -eq 2 ]; then
    # Start and end dates provided
    "$CLIENT_BIN" -start "$1" -end "$2"
elif [ $# -eq 3 ]; then
    # Start, end, and server provided
    "$CLIENT_BIN" -start "$1" -end "$2" -server "$3"
else
    echo "Usage: $0 [start_date end_date [server_address]]"
    echo ""
    echo "Examples:"
    echo "  $0                                    # Use default (last 30 days)"
    echo "  $0 2025-01-01 2025-01-31              # Query specific date range"
    echo "  $0 2025-01-01 2025-01-31 localhost:50051  # With custom server"
    echo ""
    echo "Date format: YYYY-MM-DD"
    exit 1
fi


