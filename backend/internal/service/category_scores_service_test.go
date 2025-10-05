package service

import (
	"errors"
	"testing"
	"time"

	"go-grpc-backend/internal/models"
	"go-grpc-backend/proto"
)

func TestScoreService_GetAggregatedCategoryScores_DailyGranularity(t *testing.T) {
	// Setup test dates (20 days apart - should use daily granularity)
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 1, 21, 0, 0, 0, 0, time.UTC)

	mockRepo := &mockAnalyticsRepository{
		dailyRatings: []models.CategoryRatingOverTimePeriod{
			{
				CategoryID:     1,
				CategoryName:   "Service",
				AvgPercent:     80.0,
				CategoryWeight: 0.4,
				RatingCount:    10,
				Date:           time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				CategoryID:     1,
				CategoryName:   "Service",
				AvgPercent:     90.0,
				CategoryWeight: 0.4,
				RatingCount:    15,
				Date:           time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
			},
			{
				CategoryID:     2,
				CategoryName:   "Quality",
				AvgPercent:     75.0,
				CategoryWeight: 0.6,
				RatingCount:    20,
				Date:           time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	service := newTestScoreService(mockRepo)
	result, err := service.GetAggregatedCategoryScores(startDate, endDate)

	if err != nil {
		t.Fatalf("GetAggregatedCategoryScores() error = %v", err)
	}

	// Should use daily granularity
	if result.Granularity != proto.Granularity_GRANULARITY_DAY {
		t.Errorf("Expected GRANULARITY_DAY, got %v", result.Granularity)
	}

	// Check bucket range
	if !result.BucketRange.Start.AsTime().Equal(startDate) {
		t.Errorf("Expected start date %v, got %v", startDate, result.BucketRange.Start.AsTime())
	}
	if !result.BucketRange.End.AsTime().Equal(endDate) {
		t.Errorf("Expected end date %v, got %v", endDate, result.BucketRange.End.AsTime())
	}

	// Should have 2 categories
	if len(result.Categories) != 2 {
		t.Fatalf("Expected 2 categories, got %d", len(result.Categories))
	}

	// Find Service category
	var serviceCategory *proto.CategorySeries
	for _, cat := range result.Categories {
		if cat.CategoryName == "Service" {
			serviceCategory = cat
			break
		}
	}

	if serviceCategory == nil {
		t.Fatal("Service category not found")
	}

	// Service should have 2 score points
	if len(serviceCategory.Scores) != 2 {
		t.Errorf("Expected 2 score points for Service, got %d", len(serviceCategory.Scores))
	}

	// Check that scores are sorted by date
	if len(serviceCategory.Scores) >= 2 {
		date1 := serviceCategory.Scores[0].Date.AsTime()
		date2 := serviceCategory.Scores[1].Date.AsTime()
		if !date1.Before(date2) && !date1.Equal(date2) {
			t.Errorf("Scores not sorted by date: %v should be before %v", date1, date2)
		}
	}

	// Check category total count
	if serviceCategory.CategoryTotalCount != 25 {
		t.Errorf("Expected CategoryTotalCount 25, got %d", serviceCategory.CategoryTotalCount)
	}
}

func TestScoreService_GetAggregatedCategoryScores_WeeklyGranularity(t *testing.T) {
	// Setup test dates (45 days apart - should use weekly granularity)
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 2, 15, 0, 0, 0, 0, time.UTC)

	mockRepo := &mockAnalyticsRepository{
		weeklyRatings: []models.CategoryRatingOverTimePeriod{
			{
				CategoryID:     1,
				CategoryName:   "Service",
				AvgPercent:     85.0,
				CategoryWeight: 0.5,
				RatingCount:    50,
				Date:           time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC),
			},
			{
				CategoryID:     1,
				CategoryName:   "Service",
				AvgPercent:     88.0,
				CategoryWeight: 0.5,
				RatingCount:    60,
				Date:           time.Date(2025, 1, 13, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	service := newTestScoreService(mockRepo)
	result, err := service.GetAggregatedCategoryScores(startDate, endDate)

	if err != nil {
		t.Fatalf("GetAggregatedCategoryScores() error = %v", err)
	}

	// Should use weekly granularity
	if result.Granularity != proto.Granularity_GRANULARITY_WEEK {
		t.Errorf("Expected GRANULARITY_WEEK, got %v", result.Granularity)
	}

	// Should have 1 category
	if len(result.Categories) != 1 {
		t.Fatalf("Expected 1 category, got %d", len(result.Categories))
	}

	category := result.Categories[0]

	// Check category details
	if category.CategoryId != 1 {
		t.Errorf("Expected CategoryId 1, got %d", category.CategoryId)
	}

	if category.CategoryName != "Service" {
		t.Errorf("Expected CategoryName 'Service', got %s", category.CategoryName)
	}

	// Should have 2 score points
	if len(category.Scores) != 2 {
		t.Errorf("Expected 2 score points, got %d", len(category.Scores))
	}

	// Check total count
	if category.CategoryTotalCount != 110 {
		t.Errorf("Expected CategoryTotalCount 110, got %d", category.CategoryTotalCount)
	}
}

func TestScoreService_GetAggregatedCategoryScores_MultipleCategories(t *testing.T) {
	// Setup test with multiple categories
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)

	mockRepo := &mockAnalyticsRepository{
		dailyRatings: []models.CategoryRatingOverTimePeriod{
			{
				CategoryID:     1,
				CategoryName:   "Service",
				AvgPercent:     80.0,
				CategoryWeight: 0.3,
				RatingCount:    10,
				Date:           time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				CategoryID:     2,
				CategoryName:   "Quality",
				AvgPercent:     85.0,
				CategoryWeight: 0.4,
				RatingCount:    15,
				Date:           time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				CategoryID:     3,
				CategoryName:   "Speed",
				AvgPercent:     90.0,
				CategoryWeight: 0.3,
				RatingCount:    20,
				Date:           time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	service := newTestScoreService(mockRepo)
	result, err := service.GetAggregatedCategoryScores(startDate, endDate)

	if err != nil {
		t.Fatalf("GetAggregatedCategoryScores() error = %v", err)
	}

	// Should have 3 categories
	if len(result.Categories) != 3 {
		t.Errorf("Expected 3 categories, got %d", len(result.Categories))
	}

	// Verify all categories are present
	categoryNames := make(map[string]bool)
	for _, cat := range result.Categories {
		categoryNames[cat.CategoryName] = true
	}

	expectedNames := []string{"Service", "Quality", "Speed"}
	for _, name := range expectedNames {
		if !categoryNames[name] {
			t.Errorf("Expected category %s not found", name)
		}
	}
}

func TestScoreService_GetAggregatedCategoryScores_ScoreCalculation(t *testing.T) {
	// Test score calculation formula: AvgPercent * CategoryWeight / RatingCount
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)

	avgPercent := 80.0
	categoryWeight := 0.5
	ratingCount := 10

	mockRepo := &mockAnalyticsRepository{
		dailyRatings: []models.CategoryRatingOverTimePeriod{
			{
				CategoryID:     1,
				CategoryName:   "Service",
				AvgPercent:     avgPercent,
				CategoryWeight: categoryWeight,
				RatingCount:    ratingCount,
				Date:           time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	service := newTestScoreService(mockRepo)
	result, err := service.GetAggregatedCategoryScores(startDate, endDate)

	if err != nil {
		t.Fatalf("GetAggregatedCategoryScores() error = %v", err)
	}

	if len(result.Categories) != 1 {
		t.Fatalf("Expected 1 category, got %d", len(result.Categories))
	}

	if len(result.Categories[0].Scores) != 1 {
		t.Fatalf("Expected 1 score point, got %d", len(result.Categories[0].Scores))
	}

	scorePoint := result.Categories[0].Scores[0]

	// Calculate expected score
	expectedScore := float32(avgPercent * categoryWeight / float64(ratingCount))

	if scorePoint.Score != expectedScore {
		t.Errorf("Expected score %v, got %v", expectedScore, scorePoint.Score)
	}

	// Check count
	if scorePoint.Count.GetValue() != int32(ratingCount) {
		t.Errorf("Expected count %d, got %d", ratingCount, scorePoint.Count.GetValue())
	}
}

func TestScoreService_GetAggregatedCategoryScores_EmptyResults(t *testing.T) {
	// Test with no ratings
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)

	mockRepo := &mockAnalyticsRepository{
		dailyRatings: []models.CategoryRatingOverTimePeriod{},
	}

	service := newTestScoreService(mockRepo)
	result, err := service.GetAggregatedCategoryScores(startDate, endDate)

	if err != nil {
		t.Fatalf("GetAggregatedCategoryScores() error = %v", err)
	}

	// Should return empty categories
	if len(result.Categories) != 0 {
		t.Errorf("Expected 0 categories, got %d", len(result.Categories))
	}

	// But should still have valid bucket range
	if result.BucketRange == nil {
		t.Error("Expected non-nil bucket range")
	}
}

func TestScoreService_GetAggregatedCategoryScores_DailyError(t *testing.T) {
	// Test error handling for daily granularity
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)

	expectedError := errors.New("database error")
	mockRepo := &mockAnalyticsRepository{
		dailyRatingsError: expectedError,
	}

	service := newTestScoreService(mockRepo)
	result, err := service.GetAggregatedCategoryScores(startDate, endDate)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err != expectedError {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}
}

func TestScoreService_GetAggregatedCategoryScores_WeeklyError(t *testing.T) {
	// Test error handling for weekly granularity (45 days apart)
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 2, 15, 0, 0, 0, 0, time.UTC)

	expectedError := errors.New("database error")
	mockRepo := &mockAnalyticsRepository{
		weeklyRatingsError: expectedError,
	}

	service := newTestScoreService(mockRepo)
	result, err := service.GetAggregatedCategoryScores(startDate, endDate)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err != expectedError {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}
}

func TestScoreService_GetAggregatedCategoryScores_SortingByDate(t *testing.T) {
	// Test that score points are sorted by date in ascending order
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)

	// Create data with intentionally unsorted dates
	mockRepo := &mockAnalyticsRepository{
		dailyRatings: []models.CategoryRatingOverTimePeriod{
			{
				CategoryID:     1,
				CategoryName:   "Service",
				AvgPercent:     80.0,
				CategoryWeight: 0.5,
				RatingCount:    10,
				Date:           time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC), // Middle
			},
			{
				CategoryID:     1,
				CategoryName:   "Service",
				AvgPercent:     85.0,
				CategoryWeight: 0.5,
				RatingCount:    15,
				Date:           time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), // First
			},
			{
				CategoryID:     1,
				CategoryName:   "Service",
				AvgPercent:     90.0,
				CategoryWeight: 0.5,
				RatingCount:    20,
				Date:           time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC), // Last
			},
		},
	}

	service := newTestScoreService(mockRepo)
	result, err := service.GetAggregatedCategoryScores(startDate, endDate)

	if err != nil {
		t.Fatalf("GetAggregatedCategoryScores() error = %v", err)
	}

	if len(result.Categories) != 1 {
		t.Fatalf("Expected 1 category, got %d", len(result.Categories))
	}

	scores := result.Categories[0].Scores
	if len(scores) != 3 {
		t.Fatalf("Expected 3 score points, got %d", len(scores))
	}

	// Verify scores are sorted in ascending order by date
	for i := 1; i < len(scores); i++ {
		prevDate := scores[i-1].Date.AsTime()
		currDate := scores[i].Date.AsTime()
		if !prevDate.Before(currDate) && !prevDate.Equal(currDate) {
			t.Errorf("Scores not sorted: index %d date %v should be before index %d date %v",
				i-1, prevDate, i, currDate)
		}
	}

	// Verify the specific order
	expectedDates := []time.Time{
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
	}

	for i, expectedDate := range expectedDates {
		actualDate := scores[i].Date.AsTime()
		if !actualDate.Equal(expectedDate) {
			t.Errorf("Score %d: expected date %v, got %v", i, expectedDate, actualDate)
		}
	}
}

func TestScoreService_GetAggregatedCategoryScores_ThresholdBoundary(t *testing.T) {
	// Test exactly 30 days (should use daily) and 31 days (should use weekly)
	baseDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name                string
		duration            time.Duration
		expectedGranularity proto.Granularity
	}{
		{
			name:                "Exactly 30 days - should use daily",
			duration:            30 * 24 * time.Hour,
			expectedGranularity: proto.Granularity_GRANULARITY_DAY,
		},
		{
			name:                "31 days - should use weekly",
			duration:            31 * 24 * time.Hour,
			expectedGranularity: proto.Granularity_GRANULARITY_WEEK,
		},
		{
			name:                "1 hour - should use daily",
			duration:            1 * time.Hour,
			expectedGranularity: proto.Granularity_GRANULARITY_DAY,
		},
		{
			name:                "90 days - should use weekly",
			duration:            90 * 24 * time.Hour,
			expectedGranularity: proto.Granularity_GRANULARITY_WEEK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startDate := baseDate
			endDate := baseDate.Add(tt.duration)

			mockRepo := &mockAnalyticsRepository{
				dailyRatings:  []models.CategoryRatingOverTimePeriod{},
				weeklyRatings: []models.CategoryRatingOverTimePeriod{},
			}

			service := newTestScoreService(mockRepo)
			result, err := service.GetAggregatedCategoryScores(startDate, endDate)

			if err != nil {
				t.Fatalf("GetAggregatedCategoryScores() error = %v", err)
			}

			if result.Granularity != tt.expectedGranularity {
				t.Errorf("Expected %v, got %v", tt.expectedGranularity, result.Granularity)
			}
		})
	}
}
