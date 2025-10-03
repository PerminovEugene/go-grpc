package models

import (
	"time"
)

type CategoryScore struct {
	CategoryID   int       `json:"category_id" db:"category_id"`
	CategoryName string    `json:"category_name" db:"category_name"`
	Score        float64   `json:"score" db:"score"`
	RatingCount  int       `json:"rating_count" db:"rating_count"`
	Date         time.Time `json:"date" db:"date"`
}

type TicketCategoryScore struct {
	TicketID     int     `json:"ticket_id" db:"ticket_id"`
	CategoryID   int     `json:"category_id" db:"category_id"`
	CategoryName string  `json:"category_name" db:"category_name"`
	Score        float64 `json:"score" db:"score"`
	RatingCount  int     `json:"rating_count" db:"rating_count"`
}

type OverallQualityScore struct {
	OverallScore  float64   `json:"overall_score" db:"overall_score"`
	TotalRatings  int       `json:"total_ratings" db:"total_ratings"`
	StartDate     time.Time `json:"start_date" db:"start_date"`
	EndDate       time.Time `json:"end_date" db:"end_date"`
}

type PeriodOverPeriodChange struct {
	CurrentPeriodScore  float64   `json:"current_period_score" db:"current_period_score"`
	PreviousPeriodScore float64   `json:"previous_period_score" db:"previous_period_score"`
	ChangePercentage    float64   `json:"change_percentage" db:"change_percentage"`
	CurrentStart        time.Time `json:"current_start" db:"current_start"`
	CurrentEnd          time.Time `json:"current_end" db:"current_end"`
	PreviousStart      time.Time `json:"previous_start" db:"previous_start"`
	PreviousEnd         time.Time `json:"previous_end" db:"previous_end"`
}
