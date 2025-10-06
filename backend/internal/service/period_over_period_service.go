package service

import (
	"time"

	"go-grpc-backend/internal/repository"
	"go-grpc-backend/proto"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// GetPeriodOverPeriodChange retrieves the overall quality score for two periods and calculates the change
// Returns the current period score, previous period score, and the percentage change
// Formula: ((currentScore - previousScore) / previousScore) * 100
// Uses the same scoring algorithm as GetOverallQualityScore for consistency
func GetPeriodOverPeriodChange(
	repo repository.AnalyticsRepositoryInterface,
	currentStart, currentEnd, previousStart, previousEnd time.Time,
) (*proto.PeriodOverPeriodChangeResponse, error) {
	// Get overall quality score for current period
	currentResponse, err := GetOverallQualityScore(repo, currentStart, currentEnd)
	if err != nil {
		return nil, err
	}

	// Get overall quality score for previous period
	previousResponse, err := GetOverallQualityScore(repo, previousStart, previousEnd)
	if err != nil {
		return nil, err
	}

	// Calculate percentage change
	// Formula: ((current - previous) / previous) * 100
	var changePercentage float32
	if previousResponse.OverallScore != 0 {
		changePercentage = ((currentResponse.OverallScore - previousResponse.OverallScore) / previousResponse.OverallScore) * 100
	} else {
		// If previous score is 0, we can't calculate percentage change
		// If current score is also 0, change is 0
		// If current score is > 0, we could say it's infinite growth, but we'll just set to 0
		changePercentage = 0
	}

	// Create and return response
	resp := &proto.PeriodOverPeriodChangeResponse{
		CurrentPeriodScore:   currentResponse.OverallScore,
		PreviousPeriodScore:  previousResponse.OverallScore,
		ChangePercentage:     changePercentage,
		CurrentTotalRatings:  currentResponse.TotalRatings,
		PreviousTotalRatings: previousResponse.TotalRatings,
		CurrentStart:         timestamppb.New(currentStart),
		CurrentEnd:           timestamppb.New(currentEnd),
		PreviousStart:        timestamppb.New(previousStart),
		PreviousEnd:          timestamppb.New(previousEnd),
	}

	return resp, nil
}
