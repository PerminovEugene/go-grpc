package service

import (
	"fmt"
	"time"

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

// GetAggregatedCategoryScores retrieves and processes category scores
func (s *ScoreService) GetAggregatedCategoryScores(startDate, endDate time.Time) ([]models.CategoryScore, error) {
	duration := endDate.Sub(startDate)
	useWeekly := duration > 30*24*time.Hour
	
	if useWeekly {
		return s.analyticsRepo.GetWeeklyAggregatedCategoryScores(startDate, endDate)
	}
	return s.analyticsRepo.GetAggregatedCategoryScores(startDate, endDate)
}

// GetScoresByTicket retrieves and processes ticket-based scores
func (s *ScoreService) GetScoresByTicket(startDate, endDate time.Time) ([]models.TicketCategoryScore, error) {
	return s.analyticsRepo.GetScoresByTicket(startDate, endDate)
}

// GetOverallQualityScore retrieves and processes overall quality score
func (s *ScoreService) GetOverallQualityScore(startDate, endDate time.Time) (float64, int, error) {
	return s.analyticsRepo.GetOverallQualityScore(startDate, endDate)
}

// GetPeriodOverPeriodChange calculates period-over-period change
func (s *ScoreService) GetPeriodOverPeriodChange(currentStart, currentEnd, previousStart, previousEnd time.Time) (float64, float64, error) {
	currentScore, _, err := s.analyticsRepo.GetOverallQualityScore(currentStart, currentEnd)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get current period score: %v", err)
	}

	previousScore, _, err := s.analyticsRepo.GetOverallQualityScore(previousStart, previousEnd)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get previous period score: %v", err)
	}

	return currentScore, previousScore, nil
}

// CalculateChangePercentage calculates the percentage change between two scores
func (s *ScoreService) CalculateChangePercentage(currentScore, previousScore float64) float64 {
	if previousScore > 0 {
		return ((currentScore - previousScore) / previousScore) * 100
	}
	return 0
}

// CalculateScore performs complex score calculations
// This method can be enhanced with more sophisticated scoring algorithms
func (s *ScoreService) CalculateScore() float64 {
	// This is a placeholder for more complex score calculation logic
	// In a real implementation, this might involve:
	// - Weighted averages based on category weights
	// - Time-based decay factors
	// - User behavior analysis
	// - Machine learning models
	return 1.0
}

// GetWeightedScore calculates a weighted score based on category weights
func (s *ScoreService) GetWeightedScore(scores []models.CategoryScore, categories []models.RatingCategory) float64 {
	if len(scores) == 0 || len(categories) == 0 {
		return 0.0
	}

	var totalWeightedScore float64
	var totalWeight float64

	// Create a map for quick category weight lookup
	categoryWeights := make(map[int]float64)
	for _, category := range categories {
		categoryWeights[category.ID] = float64(category.Weight)
	}

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
