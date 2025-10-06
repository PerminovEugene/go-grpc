package service

import (
	"errors"
	"testing"
	"time"

	"go-grpc-backend/internal/models"
	"go-grpc-backend/proto"
)

// Mock repository for period over period testing
type mockPeriodOverPeriodRepository struct {
	overallScoreError      error
	currentCategoryScores  []models.CategoryScore
	previousCategoryScores []models.CategoryScore
	callCount              int
}

func (m *mockPeriodOverPeriodRepository) GetOverallQualityScore(startDate, endDate time.Time) ([]models.CategoryScore, error) {
	if m.overallScoreError != nil {
		return nil, m.overallScoreError
	}

	// Return different data based on which call this is (current vs previous)
	m.callCount++
	if m.callCount == 1 {
		return m.currentCategoryScores, nil
	}
	return m.previousCategoryScores, nil
}

// testPeriodOverPeriodService wraps the logic we want to test
type testPeriodOverPeriodService struct {
	repo *mockPeriodOverPeriodRepository
}

func (s *testPeriodOverPeriodService) GetOverallQualityScore(startDate, endDate time.Time) (*proto.OverallQualityScoreResponse, error) {
	categoryScores, err := s.repo.GetOverallQualityScore(startDate, endDate)
	if err != nil {
		return nil, err
	}

	if len(categoryScores) == 0 {
		return &proto.OverallQualityScoreResponse{
			OverallScore: 0,
			TotalRatings: 0,
		}, nil
	}

	var totalScore float64
	var totalRatings int32

	for _, cs := range categoryScores {
		categoryScore := CalculateCategoryScore(cs.Score, cs.CategoryWeight)
		totalScore += categoryScore
		totalRatings += int32(cs.RatingCount)
	}

	overallScore := totalScore / float64(len(categoryScores))

	return &proto.OverallQualityScoreResponse{
		OverallScore: float32(overallScore),
		TotalRatings: totalRatings,
	}, nil
}

func (s *testPeriodOverPeriodService) GetPeriodOverPeriodChange(
	currentStart, currentEnd, previousStart, previousEnd time.Time,
) (*proto.PeriodOverPeriodChangeResponse, error) {
	// Get overall quality score for current period
	currentResponse, err := s.GetOverallQualityScore(currentStart, currentEnd)
	if err != nil {
		return nil, err
	}

	// Get overall quality score for previous period
	previousResponse, err := s.GetOverallQualityScore(previousStart, previousEnd)
	if err != nil {
		return nil, err
	}

	// Calculate percentage change
	var changePercentage float32
	if previousResponse.OverallScore != 0 {
		changePercentage = ((currentResponse.OverallScore - previousResponse.OverallScore) / previousResponse.OverallScore) * 100
	} else {
		changePercentage = 0
	}

	resp := &proto.PeriodOverPeriodChangeResponse{
		CurrentPeriodScore:   currentResponse.OverallScore,
		PreviousPeriodScore:  previousResponse.OverallScore,
		ChangePercentage:     changePercentage,
		CurrentTotalRatings:  currentResponse.TotalRatings,
		PreviousTotalRatings: previousResponse.TotalRatings,
	}

	return resp, nil
}

func newTestPeriodOverPeriodService(repo *mockPeriodOverPeriodRepository) *testPeriodOverPeriodService {
	return &testPeriodOverPeriodService{repo: repo}
}

func TestScoreService_GetPeriodOverPeriodChange_PositiveGrowth(t *testing.T) {
	currentStart := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)
	currentEnd := time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC)
	previousStart := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	previousEnd := time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)

	// Current period: score = 50 (avg 5.0, weight 0.5)
	// Previous period: score = 40 (avg 4.0, weight 0.5)
	// Change: ((50 - 40) / 40) * 100 = 25%
	mockRepo := &mockPeriodOverPeriodRepository{
		currentCategoryScores: []models.CategoryScore{
			{CategoryID: 1, CategoryName: "Service", CategoryWeight: 0.5, Score: 5.0, RatingCount: 100},
		},
		previousCategoryScores: []models.CategoryScore{
			{CategoryID: 1, CategoryName: "Service", CategoryWeight: 0.5, Score: 4.0, RatingCount: 80},
		},
	}

	service := newTestPeriodOverPeriodService(mockRepo)
	result, err := service.GetPeriodOverPeriodChange(currentStart, currentEnd, previousStart, previousEnd)

	if err != nil {
		t.Fatalf("GetPeriodOverPeriodChange() error = %v", err)
	}

	expectedCurrentScore := float32(50.0)
	expectedPreviousScore := float32(40.0)
	expectedChange := float32(25.0)

	if result.CurrentPeriodScore != expectedCurrentScore {
		t.Errorf("Expected current score %v, got %v", expectedCurrentScore, result.CurrentPeriodScore)
	}

	if result.PreviousPeriodScore != expectedPreviousScore {
		t.Errorf("Expected previous score %v, got %v", expectedPreviousScore, result.PreviousPeriodScore)
	}

	if result.ChangePercentage != expectedChange {
		t.Errorf("Expected change percentage %v, got %v", expectedChange, result.ChangePercentage)
	}

	if result.CurrentTotalRatings != 100 {
		t.Errorf("Expected 100 current ratings, got %d", result.CurrentTotalRatings)
	}

	if result.PreviousTotalRatings != 80 {
		t.Errorf("Expected 80 previous ratings, got %d", result.PreviousTotalRatings)
	}
}

func TestScoreService_GetPeriodOverPeriodChange_NegativeGrowth(t *testing.T) {
	currentStart := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)
	currentEnd := time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC)
	previousStart := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	previousEnd := time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)

	// Current period: score = 40
	// Previous period: score = 50
	// Change: ((40 - 50) / 50) * 100 = -20%
	mockRepo := &mockPeriodOverPeriodRepository{
		currentCategoryScores: []models.CategoryScore{
			{CategoryID: 1, CategoryName: "Service", CategoryWeight: 0.5, Score: 4.0, RatingCount: 80},
		},
		previousCategoryScores: []models.CategoryScore{
			{CategoryID: 1, CategoryName: "Service", CategoryWeight: 0.5, Score: 5.0, RatingCount: 100},
		},
	}

	service := newTestPeriodOverPeriodService(mockRepo)
	result, err := service.GetPeriodOverPeriodChange(currentStart, currentEnd, previousStart, previousEnd)

	if err != nil {
		t.Fatalf("GetPeriodOverPeriodChange() error = %v", err)
	}

	expectedChange := float32(-20.0)

	if result.ChangePercentage != expectedChange {
		t.Errorf("Expected change percentage %v, got %v", expectedChange, result.ChangePercentage)
	}
}

func TestScoreService_GetPeriodOverPeriodChange_NoChange(t *testing.T) {
	currentStart := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)
	currentEnd := time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC)
	previousStart := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	previousEnd := time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)

	// Both periods: score = 45
	// Change: ((45 - 45) / 45) * 100 = 0%
	mockRepo := &mockPeriodOverPeriodRepository{
		currentCategoryScores: []models.CategoryScore{
			{CategoryID: 1, CategoryName: "Service", CategoryWeight: 0.5, Score: 4.5, RatingCount: 90},
		},
		previousCategoryScores: []models.CategoryScore{
			{CategoryID: 1, CategoryName: "Service", CategoryWeight: 0.5, Score: 4.5, RatingCount: 90},
		},
	}

	service := newTestPeriodOverPeriodService(mockRepo)
	result, err := service.GetPeriodOverPeriodChange(currentStart, currentEnd, previousStart, previousEnd)

	if err != nil {
		t.Fatalf("GetPeriodOverPeriodChange() error = %v", err)
	}

	expectedChange := float32(0.0)

	if result.ChangePercentage != expectedChange {
		t.Errorf("Expected change percentage %v, got %v", expectedChange, result.ChangePercentage)
	}
}

func TestScoreService_GetPeriodOverPeriodChange_PreviousPeriodZero(t *testing.T) {
	currentStart := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)
	currentEnd := time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC)
	previousStart := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	previousEnd := time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)

	// Current period: score = 50
	// Previous period: score = 0 (no data)
	// Change: cannot calculate percentage from 0, so should be 0
	mockRepo := &mockPeriodOverPeriodRepository{
		currentCategoryScores: []models.CategoryScore{
			{CategoryID: 1, CategoryName: "Service", CategoryWeight: 0.5, Score: 5.0, RatingCount: 100},
		},
		previousCategoryScores: []models.CategoryScore{}, // Empty = score 0
	}

	service := newTestPeriodOverPeriodService(mockRepo)
	result, err := service.GetPeriodOverPeriodChange(currentStart, currentEnd, previousStart, previousEnd)

	if err != nil {
		t.Fatalf("GetPeriodOverPeriodChange() error = %v", err)
	}

	expectedChange := float32(0.0)

	if result.ChangePercentage != expectedChange {
		t.Errorf("Expected change percentage %v when previous is 0, got %v", expectedChange, result.ChangePercentage)
	}

	if result.CurrentPeriodScore == 0 {
		t.Error("Current period score should not be 0")
	}

	if result.PreviousPeriodScore != 0 {
		t.Errorf("Expected previous period score 0, got %v", result.PreviousPeriodScore)
	}
}

func TestScoreService_GetPeriodOverPeriodChange_BothPeriodsZero(t *testing.T) {
	currentStart := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)
	currentEnd := time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC)
	previousStart := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	previousEnd := time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)

	// Both periods: no data, score = 0
	// Change: 0
	mockRepo := &mockPeriodOverPeriodRepository{
		currentCategoryScores:  []models.CategoryScore{},
		previousCategoryScores: []models.CategoryScore{},
	}

	service := newTestPeriodOverPeriodService(mockRepo)
	result, err := service.GetPeriodOverPeriodChange(currentStart, currentEnd, previousStart, previousEnd)

	if err != nil {
		t.Fatalf("GetPeriodOverPeriodChange() error = %v", err)
	}

	if result.CurrentPeriodScore != 0 {
		t.Errorf("Expected current score 0, got %v", result.CurrentPeriodScore)
	}

	if result.PreviousPeriodScore != 0 {
		t.Errorf("Expected previous score 0, got %v", result.PreviousPeriodScore)
	}

	if result.ChangePercentage != 0 {
		t.Errorf("Expected change 0, got %v", result.ChangePercentage)
	}
}

func TestScoreService_GetPeriodOverPeriodChange_MultipleCategories(t *testing.T) {
	currentStart := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)
	currentEnd := time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC)
	previousStart := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	previousEnd := time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)

	// Current period:
	// Cat1: 5.0 * 0.4 * 20 = 40
	// Cat2: 4.0 * 0.6 * 20 = 48
	// Overall: (40 + 48) / 2 = 44
	//
	// Previous period:
	// Cat1: 4.0 * 0.4 * 20 = 32
	// Cat2: 3.0 * 0.6 * 20 = 36
	// Overall: (32 + 36) / 2 = 34
	//
	// Change: ((44 - 34) / 34) * 100 = 29.411764%
	mockRepo := &mockPeriodOverPeriodRepository{
		currentCategoryScores: []models.CategoryScore{
			{CategoryID: 1, CategoryName: "Service", CategoryWeight: 0.4, Score: 5.0, RatingCount: 100},
			{CategoryID: 2, CategoryName: "Quality", CategoryWeight: 0.6, Score: 4.0, RatingCount: 150},
		},
		previousCategoryScores: []models.CategoryScore{
			{CategoryID: 1, CategoryName: "Service", CategoryWeight: 0.4, Score: 4.0, RatingCount: 80},
			{CategoryID: 2, CategoryName: "Quality", CategoryWeight: 0.6, Score: 3.0, RatingCount: 120},
		},
	}

	service := newTestPeriodOverPeriodService(mockRepo)
	result, err := service.GetPeriodOverPeriodChange(currentStart, currentEnd, previousStart, previousEnd)

	if err != nil {
		t.Fatalf("GetPeriodOverPeriodChange() error = %v", err)
	}

	expectedCurrentScore := float32(44.0)
	expectedPreviousScore := float32(34.0)
	expectedChange := ((expectedCurrentScore - expectedPreviousScore) / expectedPreviousScore) * 100

	if result.CurrentPeriodScore != expectedCurrentScore {
		t.Errorf("Expected current score %v, got %v", expectedCurrentScore, result.CurrentPeriodScore)
	}

	if result.PreviousPeriodScore != expectedPreviousScore {
		t.Errorf("Expected previous score %v, got %v", expectedPreviousScore, result.PreviousPeriodScore)
	}

	// Allow for small floating point differences
	diff := result.ChangePercentage - expectedChange
	if diff < -0.001 || diff > 0.001 {
		t.Errorf("Expected change percentage %v, got %v", expectedChange, result.ChangePercentage)
	}

	if result.CurrentTotalRatings != 250 {
		t.Errorf("Expected 250 current ratings, got %d", result.CurrentTotalRatings)
	}

	if result.PreviousTotalRatings != 200 {
		t.Errorf("Expected 200 previous ratings, got %d", result.PreviousTotalRatings)
	}
}

func TestScoreService_GetPeriodOverPeriodChange_CurrentPeriodError(t *testing.T) {
	currentStart := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)
	currentEnd := time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC)
	previousStart := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	previousEnd := time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)

	expectedError := errors.New("database error")
	mockRepo := &mockPeriodOverPeriodRepository{
		overallScoreError: expectedError,
	}

	service := newTestPeriodOverPeriodService(mockRepo)
	result, err := service.GetPeriodOverPeriodChange(currentStart, currentEnd, previousStart, previousEnd)

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

func TestScoreService_GetPeriodOverPeriodChange_ExampleUseCase_WeekOverWeek(t *testing.T) {
	// Example: Current week 96% vs Previous week 95%
	// Change: ((96 - 95) / 95) * 100 = 1.0526%
	currentStart := time.Date(2025, 2, 8, 0, 0, 0, 0, time.UTC)
	currentEnd := time.Date(2025, 2, 15, 0, 0, 0, 0, time.UTC)
	previousStart := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)
	previousEnd := time.Date(2025, 2, 8, 0, 0, 0, 0, time.UTC)

	mockRepo := &mockPeriodOverPeriodRepository{
		currentCategoryScores: []models.CategoryScore{
			{CategoryID: 1, CategoryName: "Service", CategoryWeight: 1.0, Score: 4.8, RatingCount: 200},
		},
		previousCategoryScores: []models.CategoryScore{
			{CategoryID: 1, CategoryName: "Service", CategoryWeight: 1.0, Score: 4.75, RatingCount: 190},
		},
	}

	service := newTestPeriodOverPeriodService(mockRepo)
	result, err := service.GetPeriodOverPeriodChange(currentStart, currentEnd, previousStart, previousEnd)

	if err != nil {
		t.Fatalf("GetPeriodOverPeriodChange() error = %v", err)
	}

	expectedCurrentScore := float32(96.0)  // 4.8 * 1.0 * 20
	expectedPreviousScore := float32(95.0) // 4.75 * 1.0 * 20
	expectedChange := ((expectedCurrentScore - expectedPreviousScore) / expectedPreviousScore) * 100

	if result.CurrentPeriodScore != expectedCurrentScore {
		t.Errorf("Expected current score %v, got %v", expectedCurrentScore, result.CurrentPeriodScore)
	}

	if result.PreviousPeriodScore != expectedPreviousScore {
		t.Errorf("Expected previous score %v, got %v", expectedPreviousScore, result.PreviousPeriodScore)
	}

	// Allow for small floating point differences
	diff := result.ChangePercentage - expectedChange
	if diff < -0.01 || diff > 0.01 {
		t.Errorf("Expected change percentage ~%v, got %v", expectedChange, result.ChangePercentage)
	}
}

func TestScoreService_GetPeriodOverPeriodChange_LargePositiveChange(t *testing.T) {
	// Test with 100% increase (doubling)
	currentStart := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)
	currentEnd := time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC)
	previousStart := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	previousEnd := time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)

	// Current: 40, Previous: 20, Change: 100%
	mockRepo := &mockPeriodOverPeriodRepository{
		currentCategoryScores: []models.CategoryScore{
			{CategoryID: 1, CategoryName: "Service", CategoryWeight: 0.5, Score: 4.0, RatingCount: 100},
		},
		previousCategoryScores: []models.CategoryScore{
			{CategoryID: 1, CategoryName: "Service", CategoryWeight: 0.5, Score: 2.0, RatingCount: 50},
		},
	}

	service := newTestPeriodOverPeriodService(mockRepo)
	result, err := service.GetPeriodOverPeriodChange(currentStart, currentEnd, previousStart, previousEnd)

	if err != nil {
		t.Fatalf("GetPeriodOverPeriodChange() error = %v", err)
	}

	expectedChange := float32(100.0)

	if result.ChangePercentage != expectedChange {
		t.Errorf("Expected change percentage %v, got %v", expectedChange, result.ChangePercentage)
	}
}

func TestScoreService_GetPeriodOverPeriodChange_LargeNegativeChange(t *testing.T) {
	// Test with 50% decrease
	currentStart := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)
	currentEnd := time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC)
	previousStart := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	previousEnd := time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)

	// Current: 20, Previous: 40, Change: -50%
	mockRepo := &mockPeriodOverPeriodRepository{
		currentCategoryScores: []models.CategoryScore{
			{CategoryID: 1, CategoryName: "Service", CategoryWeight: 0.5, Score: 2.0, RatingCount: 50},
		},
		previousCategoryScores: []models.CategoryScore{
			{CategoryID: 1, CategoryName: "Service", CategoryWeight: 0.5, Score: 4.0, RatingCount: 100},
		},
	}

	service := newTestPeriodOverPeriodService(mockRepo)
	result, err := service.GetPeriodOverPeriodChange(currentStart, currentEnd, previousStart, previousEnd)

	if err != nil {
		t.Fatalf("GetPeriodOverPeriodChange() error = %v", err)
	}

	expectedChange := float32(-50.0)

	if result.ChangePercentage != expectedChange {
		t.Errorf("Expected change percentage %v, got %v", expectedChange, result.ChangePercentage)
	}
}
