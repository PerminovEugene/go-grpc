package service

import (
	"time"

	"go-grpc-backend/proto"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// GetOverallQualityScore retrieves the overall aggregate score for a given period
// Calculates the average of all category scores (weighted by category weight)
// Formula: (sum of all category scores) / number of categories
// Where each category score = AvgPercent * CategoryWeight * RATING_TO_PERCENT_MODIFICATOR
func (s *ScoreService) GetOverallQualityScore(startDate, endDate time.Time) (*proto.OverallQualityScoreResponse, error) {
	// Get category-level data from repository
	categoryScores, err := s.analyticsRepo.GetOverallQualityScore(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// If no categories found, return zero score
	if len(categoryScores) == 0 {
		return &proto.OverallQualityScoreResponse{
			OverallScore: 0,
			TotalRatings: 0,
			StartDate:    timestamppb.New(startDate),
			EndDate:      timestamppb.New(endDate),
		}, nil
	}

	// Calculate weighted score for each category and sum them
	var totalScore float64
	var totalRatings int32

	for _, cs := range categoryScores {
		// Use the standard formula: AvgPercent * CategoryWeight * RATING_TO_PERCENT_MODIFICATOR
		categoryScore := CalculateCategoryScore(cs.Score, cs.CategoryWeight)
		totalScore += categoryScore
		totalRatings += int32(cs.RatingCount)
	}

	// Calculate the overall score as average of all category scores
	overallScore := totalScore / float64(len(categoryScores))

	// Create and return response
	resp := &proto.OverallQualityScoreResponse{
		OverallScore: float32(overallScore),
		TotalRatings: totalRatings,
		StartDate:    timestamppb.New(startDate),
		EndDate:      timestamppb.New(endDate),
	}

	return resp, nil
}
