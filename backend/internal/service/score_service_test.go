package service

import (
	"testing"
	"time"

	"go-grpc-backend/internal/models"
)

// Mock repository for testing
type mockAnalyticsRepository struct {
	categoryScores []models.CategoryScore
	ticketScores   []models.TicketCategoryScore
	overallScore   float64
	totalRatings   int
}

func (m *mockAnalyticsRepository) GetAggregatedCategoryScores(startDate, endDate time.Time) ([]models.CategoryScore, error) {
	return m.categoryScores, nil
}

func (m *mockAnalyticsRepository) GetWeeklyAggregatedCategoryScores(startDate, endDate time.Time) ([]models.CategoryScore, error) {
	return m.categoryScores, nil
}

func (m *mockAnalyticsRepository) GetScoresByTicket(startDate, endDate time.Time) ([]models.TicketCategoryScore, error) {
	return m.ticketScores, nil
}

func (m *mockAnalyticsRepository) GetOverallQualityScore(startDate, endDate time.Time) (float64, int, error) {
	return m.overallScore, m.totalRatings, nil
}

func (m *mockAnalyticsRepository) GetRatingCategories() ([]models.RatingCategory, error) {
	return []models.RatingCategory{
		{ID: 1, Name: "Service", Weight: 4},
		{ID: 2, Name: "Quality", Weight: 6},
	}, nil
}

func TestScoreService_CalculateChangePercentage(t *testing.T) {
	service := &ScoreService{}
	
	tests := []struct {
		name           string
		currentScore   float64
		previousScore  float64
		expectedResult float64
	}{
		{
			name:           "Normal case",
			currentScore:   8.0,
			previousScore:   6.0,
			expectedResult: 33.33333333333333,
		},
		{
			name:           "Zero previous score",
			currentScore:   5.0,
			previousScore:   0.0,
			expectedResult: 0.0,
		},
		{
			name:           "Negative change",
			currentScore:   4.0,
			previousScore:   8.0,
			expectedResult: -50.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.CalculateChangePercentage(tt.currentScore, tt.previousScore)
			// Use approximate comparison for floating point values
			if result < tt.expectedResult-0.0001 || result > tt.expectedResult+0.0001 {
				t.Errorf("CalculateChangePercentage() = %v, want %v", result, tt.expectedResult)
			}
		})
	}
}

func TestScoreService_GetWeightedScore(t *testing.T) {
	service := &ScoreService{}
	
	scores := []models.CategoryScore{
		{CategoryID: 1, CategoryName: "Service", Score: 8.0, RatingCount: 10},
		{CategoryID: 2, CategoryName: "Quality", Score: 6.0, RatingCount: 15},
	}
	
	categories := []models.RatingCategory{
		{ID: 1, Name: "Service", Weight: 4},
		{ID: 2, Name: "Quality", Weight: 6},
	}
	
	weightedScore := service.GetWeightedScore(scores, categories)
	expectedScore := (8.0*4.0 + 6.0*6.0) / (4.0 + 6.0) // 6.8
	
	if weightedScore != expectedScore {
		t.Errorf("GetWeightedScore() = %v, want %v", weightedScore, expectedScore)
	}
}

func TestScoreService_GetWeightedScore_EmptyInput(t *testing.T) {
	service := &ScoreService{}
	
	// Test with empty scores
	weightedScore := service.GetWeightedScore([]models.CategoryScore{}, []models.RatingCategory{})
	if weightedScore != 0.0 {
		t.Errorf("GetWeightedScore() with empty input = %v, want 0.0", weightedScore)
	}
}
