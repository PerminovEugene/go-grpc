# gRPC Client for Analytics Service

This directory contains console clients for the Analytics gRPC service.

## Available Clients

1. **Category Scores Client** - Fetch aggregated category scores over time
2. **Ticket Scores Client** - Fetch category scores grouped by ticket
3. **Overall Quality Score Client** - Fetch overall quality score for a period
4. **Period Over Period Client** - Compare quality scores between two periods

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

---

## Overall Quality Score Client

The `overall_quality_score_client` allows you to fetch the overall quality score for a given period.

### Building

```bash
# From the backend directory
make build-overall-quality-client

# Or manually
go build -o bin/overall_quality_score_client ./client/overall_quality_score
```

### Usage

```bash
# Basic usage with date range
./bin/overall_quality_score_client -start 2025-01-01 -end 2025-01-31

# With custom server address
./bin/overall_quality_score_client -start 2025-01-01 -end 2025-01-31 -server localhost:50051

# Use default dates (last 7 days)
./bin/overall_quality_score_client
```

### Flags

- `-server`: gRPC server address (default: `localhost:50051`)
- `-start`: Start date in `YYYY-MM-DD` format (required, unless using default)
- `-end`: End date in `YYYY-MM-DD` format (required, unless using default)

### Examples

**Query past week:**
```bash
./bin/overall_quality_score_client -start 2025-01-01 -end 2025-01-07
```

**Use default (last 7 days):**
```bash
./bin/overall_quality_score_client
```

### Output Format

The client displays:
- Period information and duration
- Overall quality score (0-100%)
- Total number of ratings
- Average rating (0-5.0)
- Performance assessment

### Example Output

```
Requesting overall quality score from 2025-01-01 to 2025-01-07...

=== Overall Quality Score ===
Period: 2025-01-01 to 2025-01-07
Duration: past week

Overall Quality Score: 84.50%
Total ratings: 250
Average rating: 4.23 / 5.0

Performance: Excellent

✅ Request completed successfully
```

---

## Period Over Period Client

The `period_over_period_client` allows you to compare quality scores between two time periods and see the percentage change.

### Building

```bash
# From the backend directory
make build-period-over-period-client

# Or manually
go build -o bin/period_over_period_client ./client/period_over_period
```

### Usage

```bash
# Basic usage with both periods specified
./bin/period_over_period_client \
  -current-start 2025-02-01 -current-end 2025-02-28 \
  -previous-start 2025-01-01 -previous-end 2025-01-31

# With custom server address
./bin/period_over_period_client \
  -current-start 2025-02-01 -current-end 2025-02-28 \
  -previous-start 2025-01-01 -previous-end 2025-01-31 \
  -server localhost:50051

# Use default dates (current week vs previous week)
./bin/period_over_period_client
```

### Flags

- `-server`: gRPC server address (default: `localhost:50051`)
- `-current-start`: Current period start date in `YYYY-MM-DD` format
- `-current-end`: Current period end date in `YYYY-MM-DD` format
- `-previous-start`: Previous period start date in `YYYY-MM-DD` format
- `-previous-end`: Previous period end date in `YYYY-MM-DD` format

**Note:** All four dates must be provided together, or none for default.

### Examples

**Compare February vs January:**
```bash
./bin/period_over_period_client \
  -current-start 2025-02-01 -current-end 2025-02-28 \
  -previous-start 2025-01-01 -previous-end 2025-01-31
```

**Compare this week vs last week (using Makefile):**
```bash
make run-period-over-period-client \
  CURR_START=2025-02-08 CURR_END=2025-02-15 \
  PREV_START=2025-02-01 PREV_END=2025-02-08
```

**Use default (current week vs previous week):**
```bash
./bin/period_over_period_client
```

### Output Format

The client displays:
- Current period information (dates, score, ratings, average)
- Previous period information (dates, score, ratings, average)
- Period over period change analysis:
  - Score difference in points
  - Percentage change
  - Rating count change
  - Trend assessment (improvement/decline indicators)

### Example Output

```
Requesting period over period change...
Current period:  2025-02-01 to 2025-02-28
Previous period: 2025-01-01 to 2025-01-31

=== Period Over Period Score Change ===

📊 Current Period
   Period: 2025-02-01 to 2025-02-28
   Duration: ~1 month (28 days)
   Overall Score: 86.50%
   Total Ratings: 320
   Average Rating: 4.33 / 5.0

📊 Previous Period
   Period: 2025-01-01 to 2025-01-31
   Duration: ~1 month (31 days)
   Overall Score: 82.00%
   Total Ratings: 285
   Average Rating: 4.10 / 5.0

📈 Period Over Period Change
   Score Difference: +4.50 points
   Percentage Change: ↑ +5.49%
   Rating Count Change: +35 ratings

   Trend: ⬆️  Moderate Improvement

✅ Request completed successfully
```

### Trend Indicators

The client provides visual trend indicators:
- 🚀 Significant Improvement (≥20% increase)
- 📈 Strong Improvement (≥10% increase)
- ⬆️ Moderate Improvement (≥5% increase)
- ↗️ Slight Improvement (>0% increase)
- ➡️ Stable (0% change)
- ↘️ Slight Decline (>-5% decrease)
- ⬇️ Moderate Decline (>-10% decrease)
- 📉 Strong Decline (>-20% decrease)
- ⚠️ Significant Decline (<-20% decrease)

