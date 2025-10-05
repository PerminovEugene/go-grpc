package service

import (
	"go-grpc-backend/internal/models"
)

type CategoryIdToWeight = map[int]float64

func BuildCategoryIdToWeightMap(categories []models.RatingCategory) CategoryIdToWeight {
	categoryWeights := make(CategoryIdToWeight)
	for _, category := range categories {
		categoryWeights[category.ID] = float64(category.Weight)
	}
	return categoryWeights
}

// GetWeightedScore calculates a weighted score based on category weights
func (s *ScoreService) CalculateScore(ratings []models.CategoryScore, categoryMap CategoryIdToWeight) float64 {
	if len(ratings) == 0 {
		return 0.0
	}

	var totalWeightedScore float64
	var totalWeight float64

	for _, score := range scores {
		if weight, exists := categoryWeights[score.CategoryID]; exists {
			totalWeightedScore += score.Score * weight
			totalWeight += weight
		}
	}

	if totalWeight > 0 {
		return totalWeightedScore / totalWeight
	}
	return 0.0
}
