# gRPC Client for Analytics Service

This directory contains console clients for the Analytics gRPC service.

## Available Clients

1. **Category Scores Client** - Fetch aggregated category scores over time
2. **Ticket Scores Client** - Fetch category scores grouped by ticket

## Category Scores Client

The `category_scores_client` allows you to fetch aggregated category scores from the command line.

### Building

```bash
# From the backend directory
make build-client

# Or manually
go build -o bin/category_scores_client ./client/category_scores
```

### Usage

```bash
# Basic usage with date range
./bin/category_scores_client -start 2025-01-01 -end 2025-01-31

# With custom server address
./bin/category_scores_client -start 2025-01-01 -end 2025-03-01 -server localhost:50051

# Use default dates (last 30 days)
./bin/category_scores_client
```

### Flags

- `-server`: gRPC server address (default: `localhost:50051`)
- `-start`: Start date in `YYYY-MM-DD` format (required, unless using default)
- `-end`: End date in `YYYY-MM-DD` format (required, unless using default)

### Examples

**Query last 30 days:**
```bash
./bin/category_scores_client -start 2025-01-01 -end 2025-01-31
```

**Query 3 months (will use weekly granularity):**
```bash
./bin/category_scores_client -start 2025-01-01 -end 2025-04-01
```

**Connect to different server:**
```bash
./bin/category_scores_client -start 2025-01-01 -end 2025-01-31 -server prod-server:50051
```

### Output Format

The client displays:
- Period information and granularity (daily or weekly)
- Total number of data points
- For each category:
  - Number of data points
  - Average, minimum, and maximum scores
  - Total number of ratings
  - Recent data points (last 5)

### Example Output

```
Requesting aggregated category scores from 2025-01-01 to 2025-01-31...

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
     2025-01-29: Score=4.20, Ratings=14
     2025-01-30: Score=4.35, Ratings=15
     2025-01-31: Score=4.40, Ratings=17

📊 Category: Quality
   Number of data points: 31
   Average Score: 4.15
   Min Score: 3.60
   Max Score: 4.80
   Total Ratings: 430
   Recent data points:
     2025-01-27: Score=4.20, Ratings=14
     2025-01-28: Score=4.25, Ratings=15
     2025-01-29: Score=4.10, Ratings=13
     2025-01-30: Score=4.30, Ratings=16
     2025-01-31: Score=4.35, Ratings=18

✅ Request completed successfully
```

---

## Ticket Scores Client

The `ticket_scores_client` allows you to fetch category scores grouped by ticket from the command line.

### Building

```bash
# From the backend directory
make build-ticket-client

# Or manually
go build -o bin/ticket_scores_client ./client/ticket_scores
```

### Usage

```bash
# Basic usage with date range
./bin/ticket_scores_client -start 2025-01-01 -end 2025-01-31

# With custom server address
./bin/ticket_scores_client -start 2025-01-01 -end 2025-03-01 -server localhost:50051

# Use default dates (last 30 days)
./bin/ticket_scores_client
```

### Flags

- `-server`: gRPC server address (default: `localhost:50051`)
- `-start`: Start date in `YYYY-MM-DD` format (required, unless using default)
- `-end`: End date in `YYYY-MM-DD` format (required, unless using default)

### Examples

**Query last 30 days:**
```bash
./bin/ticket_scores_client -start 2025-01-01 -end 2025-01-31
```

**Connect to different server:**
```bash
./bin/ticket_scores_client -start 2025-01-01 -end 2025-01-31 -server prod-server:50051
```

### Output Format

The client displays:
- Period information
- Total number of tickets
- For each ticket:
  - Ticket ID
  - Category scores (as percentages 0-100%)
  - Number of ratings per category

### Example Output

```
Requesting scores by ticket from 2025-01-01 to 2025-01-31...

=== Scores by Ticket ===
Period: 2025-01-01 to 2025-01-31
Total tickets: 5

Ticket ID  | Category                  | Score      | Ratings
-----------|---------------------------|------------|------------
1          | Grammar                   |     80.00% |          4
           | Problem Solving           |     60.00% |          3
           | Tone                      |    100.00% |          5
           |                           |            |
2          | Grammar                   |     60.00% |          3
           | Problem Solving           |     40.00% |          2
           |                           |            |
3          | Grammar                   |     80.00% |          4
           | Tone                      |    100.00% |          5

✅ Request completed successfully
```

