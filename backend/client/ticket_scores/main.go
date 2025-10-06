package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sort"
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
	req := &proto.ScoresByTicketRequest{
		StartDate: timestamppb.New(start),
		EndDate:   timestamppb.New(end),
	}

	// Call the service
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Printf("Requesting scores by ticket from %s to %s...\n\n", start.Format("2006-01-02"), end.Format("2006-01-02"))

	resp, err := client.GetScoresByTicket(ctx, req)
	if err != nil {
		log.Fatalf("Failed to get scores by ticket: %v", err)
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

func displayResults(resp *proto.ScoresByTicketResponse, start, end time.Time) {
	fmt.Printf("=== Scores by Ticket ===\n")
	fmt.Printf("Period: %s to %s\n", start.Format("2006-01-02"), end.Format("2006-01-02"))
	fmt.Printf("Total tickets: %d\n\n", len(resp.Tickets))

	if len(resp.Tickets) == 0 {
		fmt.Println("No tickets found for the specified period.")
		return
	}

	// Sort tickets by ticket ID for consistent display
	tickets := make([]*proto.TicketScore, len(resp.Tickets))
	copy(tickets, resp.Tickets)
	sort.Slice(tickets, func(i, j int) bool {
		return tickets[i].TicketId < tickets[j].TicketId
	})

	// Display header for table format
	fmt.Printf("%-10s | %-25s | %-10s | %-10s\n", "Ticket ID", "Category", "Score", "Ratings")
	fmt.Println("-----------|---------------------------|------------|------------")

	// Display each ticket with its category scores
	for _, ticket := range tickets {
		if len(ticket.CategoryScores) == 0 {
			fmt.Printf("%-10d | %-25s | %-10s | %-10s\n", ticket.TicketId, "N/A", "N/A", "N/A")
			continue
		}

		// Sort categories by name for consistent display
		categories := make([]*proto.CategoryScoreForTicket, len(ticket.CategoryScores))
		copy(categories, ticket.CategoryScores)
		sort.Slice(categories, func(i, j int) bool {
			return categories[i].CategoryName < categories[j].CategoryName
		})

		// Display first category on same line as ticket ID
		firstCat := categories[0]
		fmt.Printf("%-10d | %-25s | %9.2f%% | %10d\n",
			ticket.TicketId,
			firstCat.CategoryName,
			firstCat.Score,
			firstCat.RatingCount)

		// Display remaining categories indented
		for i := 1; i < len(categories); i++ {
			cat := categories[i]
			fmt.Printf("%-10s | %-25s | %9.2f%% | %10d\n",
				"",
				cat.CategoryName,
				cat.Score,
				cat.RatingCount)
		}

		// Add separator between tickets
		if ticket != tickets[len(tickets)-1] {
			fmt.Println("           |                           |            |")
		}
	}

	fmt.Printf("\n✅ Request completed successfully\n")
}
