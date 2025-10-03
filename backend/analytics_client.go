package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go-grpc-backend/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	fmt.Println("Analytics gRPC Client Test")
	fmt.Println("==========================")

	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewAnalyticsServiceClient(conn)
	ctx := context.Background()

	startDate := time.Now().AddDate(0, 0, -7)
	endDate := time.Now()

	fmt.Printf("\nTesting Analytics Methods (Period: %v to %v)\n", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	// Test GetAggregatedCategoryScores
	fmt.Println("\n1. GetAggregatedCategoryScores")
	aggResp, err := client.GetAggregatedCategoryScores(ctx, &proto.AggregatedCategoryScoresRequest{
		StartDate: timestamppb.New(startDate),
		EndDate:   timestamppb.New(endDate),
	})
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Found %d category scores\n", len(aggResp.Scores))
		for _, score := range aggResp.Scores {
			fmt.Printf("   %s: %.1f%% (%d ratings) on %v\n", 
				score.CategoryName, score.Score, score.RatingCount, score.Date.AsTime().Format("2006-01-02"))
		}
	}

	// Test GetScoresByTicket
	fmt.Println("\n2. GetScoresByTicket")
	ticketResp, err := client.GetScoresByTicket(ctx, &proto.ScoresByTicketRequest{
		StartDate: timestamppb.New(startDate),
		EndDate:   timestamppb.New(endDate),
	})
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Found %d ticket scores\n", len(ticketResp.Scores))
		for _, score := range ticketResp.Scores {
			fmt.Printf("   Ticket %d - %s: %.1f%% (%d ratings)\n", 
				score.TicketId, score.CategoryName, score.Score, score.RatingCount)
		}
	}

	// Test GetOverallQualityScore
	fmt.Println("\n3. GetOverallQualityScore")
	overallResp, err := client.GetOverallQualityScore(ctx, &proto.OverallQualityScoreRequest{
		StartDate: timestamppb.New(startDate),
		EndDate:   timestamppb.New(endDate),
	})
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Overall Score: %.1f%% (from %d total ratings)\n", 
			overallResp.OverallScore, overallResp.TotalRatings)
	}

	// Test GetPeriodOverPeriodChange
	fmt.Println("\n4. GetPeriodOverPeriodChange")
	currentStart := time.Now().AddDate(0, 0, -7)
	currentEnd := time.Now()
	previousStart := time.Now().AddDate(0, 0, -14)
	previousEnd := time.Now().AddDate(0, 0, -7)
	
	changeResp, err := client.GetPeriodOverPeriodChange(ctx, &proto.PeriodOverPeriodChangeRequest{
		CurrentStart:  timestamppb.New(currentStart),
		CurrentEnd:    timestamppb.New(currentEnd),
		PreviousStart: timestamppb.New(previousStart),
		PreviousEnd:   timestamppb.New(previousEnd),
	})
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Current Period: %.1f%%, Previous Period: %.1f%%, Change: %.1f%%\n", 
			changeResp.CurrentPeriodScore, changeResp.PreviousPeriodScore, changeResp.ChangePercentage)
	}

	fmt.Println("\nAnalytics client test completed!")
}