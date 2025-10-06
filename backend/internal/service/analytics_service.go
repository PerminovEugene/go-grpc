package service

import (
	"go-grpc-backend/internal/models"
	"go-grpc-backend/internal/repository"
)

type ScoreService struct {
	analyticsRepo *repository.AnalyticsRepository
}

func NewScoreService(analyticsRepo *repository.AnalyticsRepository) *ScoreService {
	return &ScoreService{
		analyticsRepo: analyticsRepo,
	}
}

// CalculateChangePercentage calculates the percentage change between current and previous scores
func (s *ScoreService) CalculateChangePercentage(currentScore, previousScore float64) float64 {
	if previousScore == 0 {
		return 0.0
	}
	return ((currentScore - previousScore) / previousScore) * 100.0
}

// GetWeightedScore calculates the weighted average score across multiple categories
func (s *ScoreService) GetWeightedScore(scores []models.CategoryScore, categories []models.RatingCategory) float64 {
	if len(scores) == 0 || len(categories) == 0 {
		return 0.0
	}

	// Create a map for quick weight lookup
	weightMap := make(map[int]float64)
	for _, cat := range categories {
		weightMap[cat.ID] = float64(cat.Weight)
	}

	var totalWeightedScore float64
	var totalWeight float64

	for _, score := range scores {
		if weight, exists := weightMap[score.CategoryID]; exists {
			totalWeightedScore += score.Score * weight
			totalWeight += weight
		}
	}

	if totalWeight == 0 {
		return 0.0
	}

	return totalWeightedScore / totalWeight
}
