# Client Usage Examples

## Building and Running

### Build the clients
```bash
cd backend

# Build category scores client
make build-client

# Build ticket scores client
make build-ticket-client
```

### Run with Makefile
```bash
# Category scores client
make run-client                              # Default dates (last 30 days)
make run-client START=2025-01-01 END=2025-01-31  # With specific dates

# Ticket scores client
make run-ticket-client                       # Default dates (last 30 days)
make run-ticket-client START=2025-01-01 END=2025-01-31  # With specific dates
```

### Run directly
```bash
# Category scores client
./bin/category_scores_client -start 2025-01-01 -end 2025-01-31
./client/run-client.sh 2025-01-01 2025-01-31

# Ticket scores client
./bin/ticket_scores_client -start 2025-01-01 -end 2025-01-31
./client/run-ticket-client.sh 2025-01-01 2025-01-31
```

## Common Scenarios

### 1. Query category scores for last month
```bash
./bin/category_scores_client -start 2025-01-01 -end 2025-01-31
```

### 2. Query category scores for last quarter (will use weekly granularity)
```bash
./bin/category_scores_client -start 2025-01-01 -end 2025-04-01
```

### 3. Query ticket scores for last month
```bash
./bin/ticket_scores_client -start 2025-01-01 -end 2025-01-31
```

### 4. Query with custom server
```bash
# Category scores
./bin/category_scores_client \
  -start 2025-01-01 \
  -end 2025-01-31 \
  -server production-server:50051

# Ticket scores
./bin/ticket_scores_client \
  -start 2025-01-01 \
  -end 2025-01-31 \
  -server production-server:50051
```

### 5. Use default dates (last 30 days)
```bash
./bin/category_scores_client
./bin/ticket_scores_client
```

### 6. Testing local development
```bash
# Terminal 1 - Start the server
make run

# Terminal 2 - Run category scores client
make run-client START=2025-01-01 END=2025-01-31

# Terminal 3 - Run ticket scores client
make run-ticket-client START=2025-01-01 END=2025-01-31
```

## Understanding the Output

### Category Scores Client
The client automatically determines the granularity:
- **Daily**: For periods ≤ 30 days
- **Weekly**: For periods > 30 days

Sample Output Structure:
```
=== Aggregated Category Scores ===
Period: 2025-01-01 to 2025-01-31 (Daily granularity)
Total data points: 62

📊 Category: Service
   Number of data points: 31
   Average Score: 4.25
   Min Score: 3.80
   Max Score: 4.70
   Total Ratings: 450
   Recent data points:
     2025-01-27: Score=4.30, Ratings=15
     2025-01-28: Score=4.45, Ratings=16
     ...
```

### Ticket Scores Client
Shows category scores grouped by ticket in a table format.

Sample Output Structure:
```
=== Scores by Ticket ===
Period: 2025-01-01 to 2025-01-31
Total tickets: 5

Ticket ID  | Category                  | Score      | Ratings
-----------|---------------------------|------------|------------
1          | Grammar                   |     80.00% |          4
           | Problem Solving           |     60.00% |          3
           | Tone                      |    100.00% |          5
```

## Troubleshooting

### Connection refused
```bash
Error: Failed to connect to server: connection refused
```
**Solution**: Make sure the gRPC server is running on the specified port (default: 50051)

### Invalid date format
```bash
Error: invalid start date format (expected YYYY-MM-DD)
```
**Solution**: Use the correct date format: `2025-01-31` (not `01/31/2025` or `31-01-2025`)

### No scores found
```
No scores found for the specified period.
```
**Solution**: The database may not have data for that time period. Try a different date range.

## Scripting Examples

### Bash loop for multiple queries
```bash
#!/bin/bash
months=("01" "02" "03")
for month in "${months[@]}"; do
  echo "=== Querying month $month ==="
  echo "Category Scores:"
  ./bin/category_scores_client -start "2025-$month-01" -end "2025-$month-28"
  echo ""
  echo "Ticket Scores:"
  ./bin/ticket_scores_client -start "2025-$month-01" -end "2025-$month-28"
  echo ""
done
```

### Save output to file
```bash
# Category scores
./bin/category_scores_client \
  -start 2025-01-01 \
  -end 2025-01-31 \
  > report-category-scores-january.txt

# Ticket scores
./bin/ticket_scores_client \
  -start 2025-01-01 \
  -end 2025-01-31 \
  > report-ticket-scores-january.txt
```

### Compare different periods
```bash
#!/bin/bash
echo "=== Q1 2025 ==="
./bin/ticket_scores_client -start 2025-01-01 -end 2025-03-31
echo ""
echo "=== Q2 2025 ==="
./bin/ticket_scores_client -start 2025-04-01 -end 2025-06-30
```

