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
	// Command line flags
	var (
		serverAddr = flag.String("server", "localhost:50051", "gRPC server address")
		startDate  = flag.String("start", "", "Start date (format: 2006-01-02)")
		endDate    = flag.String("end", "", "End date (format: 2006-01-02)")
	)
	flag.Parse()

	// Parse dates
	start, end, err := parseDates(*startDate, *endDate)
	if err != nil {
		log.Fatalf("Error parsing dates: %v\n\nUsage examples:\n"+
			"  %s -start 2025-01-01 -end 2025-01-31\n"+
			"  %s -start 2025-01-01 -end 2025-03-01 -server localhost:50051\n",
			err, flag.CommandLine.Name(), flag.CommandLine.Name())
	}

	// Connect to gRPC server
	conn, err := grpc.NewClient(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Create client
	client := proto.NewAnalyticsServiceClient(conn)

	// Create request
	req := &proto.AggregatedCategoryScoresRequest{
		StartDate: timestamppb.New(start),
		EndDate:   timestamppb.New(end),
	}

	// Call the service
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Printf("Requesting aggregated category scores from %s to %s...\n\n", start.Format("2006-01-02"), end.Format("2006-01-02"))

	resp, err := client.GetAggregatedCategoryScores(ctx, req)
	if err != nil {
		log.Fatalf("Failed to get aggregated category scores: %v", err)
	}

	// Display results
	displayResults(resp, start, end)
}

func parseDates(startStr, endStr string) (time.Time, time.Time, error) {
	const layout = "2006-01-02"

	// If no dates provided, use last 30 days as default
	if startStr == "" && endStr == "" {
		end := time.Now()
		start := end.AddDate(0, 0, -30)
		fmt.Printf("No dates provided, using default: last 30 days\n")
		return start, end, nil
	}

	if startStr == "" {
		return time.Time{}, time.Time{}, fmt.Errorf("start date is required when end date is provided")
	}
	if endStr == "" {
		return time.Time{}, time.Time{}, fmt.Errorf("end date is required when start date is provided")
	}

	start, err := time.Parse(layout, startStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid start date format (expected YYYY-MM-DD): %v", err)
	}

	end, err := time.Parse(layout, endStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid end date format (expected YYYY-MM-DD): %v", err)
	}

	if end.Before(start) {
		return time.Time{}, time.Time{}, fmt.Errorf("end date must be after start date")
	}

	return start, end, nil
}

func displayResults(resp *proto.AggregatedCategoryScoresResponse, start, end time.Time) {
	granularity := "Daily"
	if resp.Granularity == proto.Granularity_GRANULARITY_WEEK {
		granularity = "Weekly"
	}

	fmt.Printf("=== Aggregated Category Scores ===\n")
	fmt.Printf("Period: %s to %s (%s granularity)\n", start.Format("2006-01-02"), end.Format("2006-01-02"), granularity)
	fmt.Printf("Total categories: %d\n\n", len(resp.Categories))

	if len(resp.Categories) == 0 {
		fmt.Println("No scores found for the specified periodd.")
		return
	}

	// Display each category
	for _, categorySeries := range resp.Categories {
		fmt.Printf("📊 Category: %s\n", categorySeries.CategoryName)
		fmt.Printf("   Total ratings in period: %d\n", categorySeries.CategoryTotalCount)
		fmt.Printf("   Number of data points: %d\n\n", len(categorySeries.Scores))

		if len(categorySeries.Scores) == 0 {
			fmt.Printf("   No score data available\n\n")
			continue
		}

		// Display each score point with date and score
		for _, scorePoint := range categorySeries.Scores {
			date := scorePoint.Date.AsTime().Format("2006-01-02")
			countStr := ""
			if scorePoint.Count != nil {
				countStr = fmt.Sprintf(" (count: %d)", scorePoint.Count.Value)
			}
			fmt.Printf("   %s: %.2f%s\n", date, scorePoint.Score, countStr)
		}

		fmt.Println()
	}

	fmt.Printf("✅ Request completed successfully\n")
}
