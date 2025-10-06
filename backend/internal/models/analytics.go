package models

import (
	"time"
)

// internal
type CategoryRatingOverTimePeriod struct {
	CategoryID        int       `json:"category_id" db:"category_id"`
	CategoryName      string    `json:"category_name" db:"category_name"`
	AvgPercent        float64   `json:"avg_percent" db:"avg_percent"`
	CategoryWeight    float64   `json:"weight" db:"weight"`
	RatingCount       int       `json:"rating_count" db:"rating_count"`
	Date              time.Time `json:"date" db:"bucket_date"`
	RatingsTotalCount int       `json:"ratings_total_count" db:"ratings_total_count"`
}

// response
type Granularity string

const (
	GranularityDay  Granularity = "day"
	GranularityWeek Granularity = "week"
)

type ScorePoint struct {
	Date  time.Time `json:"date"`
	Score float64   `json:"score"`
	Count int       `json:"count"`
}

type CategorySeries struct {
	CategoryID         int          `json:"category_id"`
	CategoryName       string       `json:"category_name"`
	CategoryTotalCount int          `json:"category_total_count"`
	Scores             []ScorePoint `json:"scores"`
	PeriodScore        *float64     `json:"period_score,omitempty"`
}

type BucketRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type AggregatedResponse struct {
	Granularity Granularity      `json:"granularity"`
	BucketRange BucketRange      `json:"bucket_range"`
	Categories  []CategorySeries `json:"categories"`
}

type TicketCategoryScore struct {
	TicketID     int     `json:"ticket_id" db:"ticket_id"`
	CategoryID   int     `json:"category_id" db:"category_id"`
	CategoryName string  `json:"category_name" db:"category_name"`
	Score        float64 `json:"score" db:"score"`
	RatingCount  int     `json:"rating_count" db:"rating_count"`
}

type OverallQualityScore struct {
	OverallScore float64   `json:"overall_score" db:"overall_score"`
	TotalRatings int       `json:"total_ratings" db:"total_ratings"`
	StartDate    time.Time `json:"start_date" db:"start_date"`
	EndDate      time.Time `json:"end_date" db:"end_date"`
}

type PeriodOverPeriodChange struct {
	CurrentPeriodScore  float64   `json:"current_period_score" db:"current_period_score"`
	PreviousPeriodScore float64   `json:"previous_period_score" db:"previous_period_score"`
	ChangePercentage    float64   `json:"change_percentage" db:"change_percentage"`
	CurrentStart        time.Time `json:"current_start" db:"current_start"`
	CurrentEnd          time.Time `json:"current_end" db:"current_end"`
	PreviousStart       time.Time `json:"previous_start" db:"previous_start"`
	PreviousEnd         time.Time `json:"previous_end" db:"previous_end"`
}

type CategoryScore struct {
	CategoryID   int     `json:"category_id" db:"category_id"`
	CategoryName string  `json:"category_name" db:"category_name"`
	Score        float64 `json:"score" db:"score"`
	RatingCount  int     `json:"rating_count" db:"rating_count"`
}
