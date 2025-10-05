# Client Usage Examples

## Building and Running

### Build the client
```bash
cd backend
make build-client
```

### Run with Makefile
```bash
# With default dates (last 30 days)
make run-client

# With specific dates
make run-client START=2025-01-01 END=2025-01-31
```

### Run directly
```bash
# Run the binary
./bin/category_scores_client -start 2025-01-01 -end 2025-01-31

# Or use the helper script
./client/run-client.sh 2025-01-01 2025-01-31
```

## Common Scenarios

### 1. Query last month
```bash
./bin/category_scores_client -start 2025-01-01 -end 2025-01-31
```

### 2. Query last quarter (will use weekly granularity)
```bash
./bin/category_scores_client -start 2025-01-01 -end 2025-04-01
```

### 3. Query with custom server
```bash
./bin/category_scores_client \
  -start 2025-01-01 \
  -end 2025-01-31 \
  -server production-server:50051
```

### 4. Use default dates (last 30 days)
```bash
./bin/category_scores_client
```

### 5. Testing local development
```bash
# Terminal 1 - Start the server
make run

# Terminal 2 - Run the client
make run-client START=2025-01-01 END=2025-01-31
```

## Understanding the Output

The client automatically determines the granularity:
- **Daily**: For periods ≤ 30 days
- **Weekly**: For periods > 30 days

### Sample Output Structure
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
  echo "Querying month $month..."
  ./bin/category_scores_client -start "2025-$month-01" -end "2025-$month-28"
  echo ""
done
```

### Save output to file
```bash
./bin/category_scores_client \
  -start 2025-01-01 \
  -end 2025-01-31 \
  > report-january.txt
```

### JSON parsing (if output is JSON in the future)
```bash
./bin/category_scores_client \
  -start 2025-01-01 \
  -end 2025-01-31 \
  -format json | jq '.scores[] | select(.category_name == "Service")'
```

