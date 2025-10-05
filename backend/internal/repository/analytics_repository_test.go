package repository

import (
	"database/sql"
	"testing"
	"time"
)

func TestAnalyticsRepository_NewAnalyticsRepository(t *testing.T) {
	// Test repository creation
	var db *sql.DB
	repo := NewAnalyticsRepository(db)

	// Verify that the repository is created correctly
	if repo == nil {
		t.Error("NewAnalyticsRepository() should not return nil")
	}

	if repo.db != db {
		t.Error("NewAnalyticsRepository() should store the database connection")
	}
}

func TestAnalyticsRepository_GetRatingCategories_ErrorHandling(t *testing.T) {
	// Test with nil database connection
	var db *sql.DB
	repo := NewAnalyticsRepository(db)

	// This will panic with nil database, so we expect a panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("GetRatingCategories() with nil database should panic")
		}
	}()

	repo.GetRatingCategories()
}

func TestAnalyticsRepository_GetOverallQualityScore_ErrorHandling(t *testing.T) {
	// Test with nil database connection
	var db *sql.DB
	repo := NewAnalyticsRepository(db)

	startDate := time.Now().AddDate(0, 0, -7)
	endDate := time.Now()

	// This will panic with nil database, so we expect a panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("GetOverallQualityScore() with nil database should panic")
		}
	}()

	repo.GetOverallQualityScore(startDate, endDate)
}
