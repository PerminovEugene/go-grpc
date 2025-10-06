package service

import (
	"time"

	"go-grpc-backend/proto"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// GetScoresByTicket retrieves and aggregates category scores by ticket for a given period
func (s *ScoreService) GetScoresByTicket(startDate, endDate time.Time) (*proto.ScoresByTicketResponse, error) {
	// Get data from repository
	scores, err := s.analyticsRepo.GetScoresByTicket(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Group scores by ticket ID
	ticketMap := make(map[int32]*proto.TicketScore)
	for _, score := range scores {
		ticketID := int32(score.TicketID)

		// Create ticket entry if it doesn't exist
		ticket, ok := ticketMap[ticketID]
		if !ok {
			ticket = &proto.TicketScore{
				TicketId:       ticketID,
				CategoryScores: nil,
			}
			ticketMap[ticketID] = ticket
		}

		// Convert avg rating (0-5) to percentage score (0-100)
		// Using the same formula as in GetAggregatedCategoryScores
		percentScore := float32(score.Score * RATING_TO_PERCENT_MODIFICATOR)

		// Add category score to ticket
		categoryScore := &proto.CategoryScoreForTicket{
			CategoryId:   int32(score.CategoryID),
			CategoryName: score.CategoryName,
			Score:        percentScore,
			RatingCount:  int32(score.RatingCount),
		}
		ticket.CategoryScores = append(ticket.CategoryScores, categoryScore)
	}

	// Convert map to slice
	tickets := make([]*proto.TicketScore, 0, len(ticketMap))
	for _, ticket := range ticketMap {
		tickets = append(tickets, ticket)
	}

	// Create and return response
	resp := &proto.ScoresByTicketResponse{
		Tickets:   tickets,
		StartDate: timestamppb.New(startDate),
		EndDate:   timestamppb.New(endDate),
	}

	return resp, nil
}
