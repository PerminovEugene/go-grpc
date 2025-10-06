package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"go-grpc-backend/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	var (
		serverAddr    = flag.String("server", "localhost:50051", "gRPC server address")
		currentStart  = flag.String("current-start", "", "Current period start date (format: 2006-01-02)")
		currentEnd    = flag.String("current-end", "", "Current period end date (format: 2006-01-02)")
		previousStart = flag.String("previous-start", "", "Previous period start date (format: 2006-01-02)")
		previousEnd   = flag.String("previous-end", "", "Previous period end date (format: 2006-01-02)")
	)
	flag.Parse()

	currStart, currEnd, prevStart, prevEnd, err := parseDates(*currentStart, *currentEnd, *previousStart, *previousEnd)
	if err != nil {
		log.Fatalf("Error parsing dates: %v\n", err)
	}

	conn, err := grpc.NewClient(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := proto.NewAnalyticsServiceClient(conn)

	req := &proto.PeriodOverPeriodChangeRequest{
		CurrentStart:  timestamppb.New(currStart),
		CurrentEnd:    timestamppb.New(currEnd),
		PreviousStart: timestamppb.New(prevStart),
		PreviousEnd:   timestamppb.New(prevEnd),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Printf("Requesting period over period change...\n")
	fmt.Printf("Current period:  %s to %s\n", currStart.Format("2006-01-02"), currEnd.Format("2006-01-02"))
	fmt.Printf("Previous period: %s to %s\n\n", prevStart.Format("2006-01-02"), prevEnd.Format("2006-01-02"))

	resp, err := client.GetPeriodOverPeriodChange(ctx, req)
	if err != nil {
		log.Fatalf("Failed to get period over period change: %v", err)
	}

	displayResults(resp, currStart, currEnd, prevStart, prevEnd)
}

func parseDates(currentStart, currentEnd, previousStart, previousEnd string) (time.Time, time.Time, time.Time, time.Time, error) {
	const layout = "2006-01-02"

	// If no dates provided, use default: current week vs previous week
	if currentStart == "" && currentEnd == "" && previousStart == "" && previousEnd == "" {
		now := time.Now()
		// Current week: last 7 days
		currEnd := now
		currStart := now.AddDate(0, 0, -7)
		// Previous week: 7 days before that
		prevEnd := currStart
		prevStart := currStart.AddDate(0, 0, -7)

		fmt.Printf("No dates provided, using default:\n")
		fmt.Printf("  Current week:  last 7 days\n")
		fmt.Printf("  Previous week: 7 days before that\n\n")

		return currStart, currEnd, prevStart, prevEnd, nil
	}

	// All four dates must be provided
	if currentStart == "" || currentEnd == "" || previousStart == "" || previousEnd == "" {
		return time.Time{}, time.Time{}, time.Time{}, time.Time{}, fmt.Errorf("all four dates required: current-start, current-end, previous-start, previous-end")
	}

	currStart, err := time.Parse(layout, currentStart)
	if err != nil {
		return time.Time{}, time.Time{}, time.Time{}, time.Time{}, fmt.Errorf("invalid current-start date: %v", err)
	}

	currEnd, err := time.Parse(layout, currentEnd)
	if err != nil {
		return time.Time{}, time.Time{}, time.Time{}, time.Time{}, fmt.Errorf("invalid current-end date: %v", err)
	}

	prevStart, err := time.Parse(layout, previousStart)
	if err != nil {
		return time.Time{}, time.Time{}, time.Time{}, time.Time{}, fmt.Errorf("invalid previous-start date: %v", err)
	}

	prevEnd, err := time.Parse(layout, previousEnd)
	if err != nil {
		return time.Time{}, time.Time{}, time.Time{}, time.Time{}, fmt.Errorf("invalid previous-end date: %v", err)
	}

	// Validate date ranges
	if currEnd.Before(currStart) {
		return time.Time{}, time.Time{}, time.Time{}, time.Time{}, fmt.Errorf("current end date must be after current start date")
	}

	if prevEnd.Before(prevStart) {
		return time.Time{}, time.Time{}, time.Time{}, time.Time{}, fmt.Errorf("previous end date must be after previous start date")
	}

	return currStart, currEnd, prevStart, prevEnd, nil
}

func displayResults(resp *proto.PeriodOverPeriodChangeResponse, currStart, currEnd, prevStart, prevEnd time.Time) {
	fmt.Printf("=== Period Over Period Score Change ===\n\n")

	// Current period
	fmt.Printf("📊 Current Period\n")
	fmt.Printf("   Period: %s to %s\n", currStart.Format("2006-01-02"), currEnd.Format("2006-01-02"))
	currDuration := currEnd.Sub(currStart)
	currDays := int(currDuration.Hours() / 24)
	fmt.Printf("   Duration: %s\n", formatDuration(currDays))
	fmt.Printf("   Overall Score: %.2f%%\n", resp.CurrentPeriodScore)
	fmt.Printf("   Total Ratings: %d\n", resp.CurrentTotalRatings)
	if resp.CurrentTotalRatings > 0 {
		fmt.Printf("   Average Rating: %.2f / 5.0\n", resp.CurrentPeriodScore/20.0)
	}
	fmt.Println()

	// Previous period
	fmt.Printf("📊 Previous Period\n")
	fmt.Printf("   Period: %s to %s\n", prevStart.Format("2006-01-02"), prevEnd.Format("2006-01-02"))
	prevDuration := prevEnd.Sub(prevStart)
	prevDays := int(prevDuration.Hours() / 24)
	fmt.Printf("   Duration: %s\n", formatDuration(prevDays))
	fmt.Printf("   Overall Score: %.2f%%\n", resp.PreviousPeriodScore)
	fmt.Printf("   Total Ratings: %d\n", resp.PreviousTotalRatings)
	if resp.PreviousTotalRatings > 0 {
		fmt.Printf("   Average Rating: %.2f / 5.0\n", resp.PreviousPeriodScore/20.0)
	}
	fmt.Println()

	// Change analysis
	fmt.Printf("📈 Period Over Period Change\n")

	// Score difference
	scoreDiff := resp.CurrentPeriodScore - resp.PreviousPeriodScore
	fmt.Printf("   Score Difference: ")
	if scoreDiff > 0 {
		fmt.Printf("+%.2f points\n", scoreDiff)
	} else if scoreDiff < 0 {
		fmt.Printf("%.2f points\n", scoreDiff)
	} else {
		fmt.Printf("No change\n")
	}

	// Percentage change
	fmt.Printf("   Percentage Change: ")
	if resp.ChangePercentage > 0 {
		fmt.Printf("↑ +%.2f%%\n", resp.ChangePercentage)
	} else if resp.ChangePercentage < 0 {
		fmt.Printf("↓ %.2f%%\n", resp.ChangePercentage)
	} else {
		fmt.Printf("→ 0.00%% (No change)\n")
	}

	// Rating count change
	ratingDiff := int(resp.CurrentTotalRatings) - int(resp.PreviousTotalRatings)
	fmt.Printf("   Rating Count Change: ")
	if ratingDiff > 0 {
		fmt.Printf("+%d ratings\n", ratingDiff)
	} else if ratingDiff < 0 {
		fmt.Printf("%d ratings\n", ratingDiff)
	} else {
		fmt.Printf("No change\n")
	}

	// Trend assessment
	fmt.Printf("\n   Trend: %s\n", getTrendAssessment(resp.ChangePercentage))

	// Special notes
	if resp.PreviousPeriodScore == 0 && resp.CurrentPeriodScore > 0 {
		fmt.Printf("\n   ⚠️  Note: Previous period had no data. Percentage change cannot be calculated.\n")
	} else if resp.CurrentTotalRatings == 0 && resp.PreviousTotalRatings == 0 {
		fmt.Printf("\n   ⚠️  Note: No ratings found in either period.\n")
	}

	fmt.Printf("\n✅ Request completed successfully\n")
}

func formatDuration(days int) string {
	switch {
	case days == 1:
		return "1 day"
	case days == 7:
		return "1 week (7 days)"
	case days >= 28 && days <= 31:
		return fmt.Sprintf("~1 month (%d days)", days)
	case days >= 365:
		years := days / 365
		if years == 1 {
			return "1 year"
		}
		return fmt.Sprintf("%d years", years)
	default:
		return fmt.Sprintf("%d days", days)
	}
}

func getTrendAssessment(changePercentage float32) string {
	switch {
	case changePercentage >= 20:
		return "🚀 Significant Improvement"
	case changePercentage >= 10:
		return "📈 Strong Improvement"
	case changePercentage >= 5:
		return "⬆️  Moderate Improvement"
	case changePercentage > 0:
		return "↗️  Slight Improvement"
	case changePercentage == 0:
		return "➡️  Stable"
	case changePercentage > -5:
		return "↘️  Slight Decline"
	case changePercentage > -10:
		return "⬇️  Moderate Decline"
	case changePercentage > -20:
		return "📉 Strong Decline"
	default:
		return "⚠️  Significant Decline"
	}
}
