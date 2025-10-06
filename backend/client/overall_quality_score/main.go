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
		serverAddr = flag.String("server", "localhost:50051", "gRPC server address")
		startDate  = flag.String("start", "", "Start date (format: 2006-01-02)")
		endDate    = flag.String("end", "", "End date (format: 2006-01-02)")
	)
	flag.Parse()

	start, end, err := parseDates(*startDate, *endDate)
	if err != nil {
		log.Fatalf("Error parsing dates: %v\n", err)
	}

	conn, err := grpc.NewClient(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := proto.NewAnalyticsServiceClient(conn)

	req := &proto.OverallQualityScoreRequest{
		StartDate: timestamppb.New(start),
		EndDate:   timestamppb.New(end),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Printf("Requesting overall quality score from %s to %s...\n\n", start.Format("2006-01-02"), end.Format("2006-01-02"))

	resp, err := client.GetOverallQualityScore(ctx, req)
	if err != nil {
		log.Fatalf("Failed to get overall quality score: %v", err)
	}

	displayResults(resp, start, end)
}

func parseDates(startStr, endStr string) (time.Time, time.Time, error) {
	const layout = "2006-01-02"

	if startStr == "" && endStr == "" {
		end := time.Now()
		start := end.AddDate(0, 0, -7)
		fmt.Printf("No dates provided, using default: last 7 days\n")
		return start, end, nil
	}

	if startStr == "" || endStr == "" {
		return time.Time{}, time.Time{}, fmt.Errorf("both start and end dates required")
	}

	start, err := time.Parse(layout, startStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid start date: %v", err)
	}

	end, err := time.Parse(layout, endStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid end date: %v", err)
	}

	if end.Before(start) {
		return time.Time{}, time.Time{}, fmt.Errorf("end date must be after start date")
	}

	return start, end, nil
}

func displayResults(resp *proto.OverallQualityScoreResponse, start, end time.Time) {
	fmt.Printf("=== Overall Quality Score ===\n")
	fmt.Printf("Period: %s to %s\n", start.Format("2006-01-02"), end.Format("2006-01-02"))

	duration := end.Sub(start)
	daysCount := int(duration.Hours() / 24)

	var periodDesc string
	switch {
	case daysCount == 7:
		periodDesc = "past week"
	case daysCount >= 28 && daysCount <= 31:
		periodDesc = "past month"
	default:
		periodDesc = fmt.Sprintf("%d days", daysCount)
	}
	fmt.Printf("Duration: %s\n\n", periodDesc)

	if resp.TotalRatings == 0 {
		fmt.Println("No ratings found for the specified period.")
		return
	}

	fmt.Printf("Overall Quality Score: %.2f%%\n", resp.OverallScore)
	fmt.Printf("Total ratings: %d\n", resp.TotalRatings)
	fmt.Printf("Average rating: %.2f / 5.0\n\n", resp.OverallScore/20.0)

	var assessment string
	switch {
	case resp.OverallScore >= 95:
		assessment = "Exceptional"
	case resp.OverallScore >= 85:
		assessment = "Excellent"
	case resp.OverallScore >= 75:
		assessment = "Good"
	case resp.OverallScore >= 65:
		assessment = "Average"
	default:
		assessment = "Needs Improvement"
	}
	fmt.Printf("Performance: %s\n\n", assessment)
	fmt.Printf("✅ Request completed successfully\n")
}
