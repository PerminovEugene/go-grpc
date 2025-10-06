# Period Over Period Client - Examples

This document provides practical examples for using the Period Over Period gRPC client.

## Quick Start

### Default Usage (Current Week vs Previous Week)

```bash
./bin/period_over_period_client
```

This will compare the last 7 days (current week) against the 7 days before that (previous week).

## Common Use Cases

### 1. Week Over Week Comparison

Compare the current week with the previous week:

```bash
./bin/period_over_period_client \
  -current-start 2025-02-08 \
  -current-end 2025-02-15 \
  -previous-start 2025-02-01 \
  -previous-end 2025-02-08
```

**Using Makefile:**
```bash
make run-period-over-period-client \
  CURR_START=2025-02-08 CURR_END=2025-02-15 \
  PREV_START=2025-02-01 PREV_END=2025-02-08
```

### 2. Month Over Month Comparison

Compare February with January:

```bash
./bin/period_over_period_client \
  -current-start 2025-02-01 \
  -current-end 2025-02-28 \
  -previous-start 2025-01-01 \
  -previous-end 2025-01-31
```

**Using Makefile:**
```bash
make run-period-over-period-client \
  CURR_START=2025-02-01 CURR_END=2025-02-28 \
  PREV_START=2025-01-01 PREV_END=2025-01-31
```

### 3. Quarter Over Quarter Comparison

Compare Q1 2025 with Q4 2024:

```bash
./bin/period_over_period_client \
  -current-start 2025-01-01 \
  -current-end 2025-03-31 \
  -previous-start 2024-10-01 \
  -previous-end 2024-12-31
```

### 4. Year Over Year Comparison

Compare 2025 with 2024:

```bash
./bin/period_over_period_client \
  -current-start 2025-01-01 \
  -current-end 2025-12-31 \
  -previous-start 2024-01-01 \
  -previous-end 2024-12-31
```

### 5. Custom Period Comparison

Compare any two custom periods (e.g., 10 days each):

```bash
./bin/period_over_period_client \
  -current-start 2025-02-10 \
  -current-end 2025-02-20 \
  -previous-start 2025-01-20 \
  -previous-end 2025-01-30
```

## Advanced Usage

### Connect to Different Server

```bash
./bin/period_over_period_client \
  -current-start 2025-02-01 \
  -current-end 2025-02-28 \
  -previous-start 2025-01-01 \
  -previous-end 2025-01-31 \
  -server production-server:50051
```

### Holiday Period Analysis

Compare a holiday period with a non-holiday period:

```bash
# December (holiday season) vs November (regular season)
./bin/period_over_period_client \
  -current-start 2024-12-01 \
  -current-end 2024-12-31 \
  -previous-start 2024-11-01 \
  -previous-end 2024-11-30
```

## Understanding the Output

### Positive Growth Example

```
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
```

**Interpretation:**
- Quality score increased from 82% to 86.5%
- This represents a 5.49% improvement
- More ratings collected in the current period (+35)
- The trend indicator shows "Moderate Improvement"

### Negative Growth Example

```
📈 Period Over Period Change
   Score Difference: -3.20 points
   Percentage Change: ↓ -4.12%
   Rating Count Change: -15 ratings

   Trend: ↘️  Slight Decline
```

**Interpretation:**
- Quality score decreased by 3.2 percentage points
- This represents a 4.12% decline
- Fewer ratings in the current period (-15)
- The trend indicator shows "Slight Decline"

### No Change Example

```
📈 Period Over Period Change
   Score Difference: No change
   Percentage Change: → 0.00% (No change)
   Rating Count Change: +5 ratings

   Trend: ➡️  Stable
```

**Interpretation:**
- Quality score remained the same
- However, slightly more ratings were collected
- Performance is stable

## Trend Indicators Guide

| Indicator | Change Range | Meaning |
|-----------|--------------|---------|
| 🚀 | ≥ +20% | Significant Improvement |
| 📈 | ≥ +10% | Strong Improvement |
| ⬆️ | ≥ +5% | Moderate Improvement |
| ↗️ | > 0% | Slight Improvement |
| ➡️ | 0% | Stable |
| ↘️ | > -5% | Slight Decline |
| ⬇️ | > -10% | Moderate Decline |
| 📉 | > -20% | Strong Decline |
| ⚠️ | ≤ -20% | Significant Decline |

## Tips for Analysis

1. **Equal Period Lengths**: For accurate comparison, use equal-length periods (e.g., both 7 days, both 30 days)

2. **Consistent Days of Week**: When comparing weeks, try to use the same days of the week (Monday-Sunday)

3. **Seasonal Factors**: Consider seasonal variations when analyzing the results
   - Holiday periods may have different patterns
   - Beginning/end of month may affect results

4. **Sample Size**: Pay attention to the rating counts
   - Lower rating counts may make percentages less reliable
   - Large differences in rating counts may indicate data quality issues

5. **Context Matters**: Always consider external factors:
   - Were there any process changes?
   - Did team composition change?
   - Were there any technical issues?

## Automation Examples

### Shell Script for Weekly Reporting

```bash
#!/bin/bash
# weekly_report.sh

# Calculate dates for current week and previous week
CURR_END=$(date +%Y-%m-%d)
CURR_START=$(date -d "$CURR_END -7 days" +%Y-%m-%d)
PREV_END=$CURR_START
PREV_START=$(date -d "$PREV_END -7 days" +%Y-%m-%d)

echo "Generating weekly report..."
./bin/period_over_period_client \
  -current-start $CURR_START \
  -current-end $CURR_END \
  -previous-start $PREV_START \
  -previous-end $PREV_END
```

### Cron Job for Daily Reports

Add to crontab to run every day at 9 AM:

```bash
0 9 * * * cd /path/to/project && make run-period-over-period-client > /var/log/daily-quality-report.log 2>&1
```

## Troubleshooting

### No Data in One or Both Periods

If you see:
```
⚠️  Note: Previous period had no data. Percentage change cannot be calculated.
```

**Solutions:**
- Check if the date ranges are correct
- Verify that ratings exist in the database for those periods
- Ensure the server is connected to the correct database

### All Four Dates Required Error

If you see:
```
Error: all four dates required: current-start, current-end, previous-start, previous-end
```

**Solution:**
Either provide all four dates or none (for default behavior).

## Related Documentation

- [Client README](../README.md) - General client documentation
- [Overall Quality Score Client](../overall_quality_score/) - For single period analysis
- [Category Scores Client](../category_scores/) - For detailed category breakdown

