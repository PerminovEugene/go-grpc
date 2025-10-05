package repository

import (
	"database/sql"
	"fmt"
	"time"

	"go-grpc-backend/internal/models"
)

type AnalyticsRepository struct {
	db *sql.DB
}

func NewAnalyticsRepository(db *sql.DB) *AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

func (r *AnalyticsRepository) GetWeeklyAggregatedCategoryRatings(
	startDate, endDate time.Time,
) ([]models.CategoryRatingOverTimePeriod, error) {

	const query = `
		SELECT
			rc.id AS category_id,
			rc.name AS category_name,
			rc.weight as category_weight
			AVG(r.rating) * 20.0 AS avg_percent,
			COUNT(r.id) AS rating_count,
			date(r.created_at, 'weekday 1', '-7 days') AS bucket_week_start,
			COUNT(*) OVER (PARTITION BY rc.id) AS ratings_total
		FROM ratings r
		JOIN rating_categories rc ON rc.id = r.rating_category_id
		WHERE r.created_at >= ? AND r.created_at < ?
		GROUP BY rc.id, rc.name, bucket_week_start
		ORDER BY rc.name, bucket_week_start;
	`

	rows, err := r.db.Query(query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("query weekly aggregated category rating: %w", err)
	}
	defer rows.Close()

	ratings := make([]models.CategoryRatingOverTimePeriod, 0, 128)

	for rows.Next() {
		var (
			score     models.CategoryRatingOverTimePeriod
			bucketStr string
		)

		if err := rows.Scan(
			&score.CategoryID,
			&score.CategoryName,
			&score.AvgPercent, 
			&score.RatingCount,
			&bucketStr, 
			&score.RatingsTotalCount, 
		); err != nil {
			return nil, fmt.Errorf("scan weekly category score: %w", err)
		}

		t, err := time.Parse("2006-01-02", bucketStr)
		if err != nil {
			return nil, fmt.Errorf("parse week start %q: %w", bucketStr, err)
		}
		score.Date = t 

		ratings = append(ratings, score)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return ratings, nil
}

func (r *AnalyticsRepository) GetDailyAggregatedCategoryRatings(startDate, endDate time.Time) ([]models.CategoryRatingOverTimePeriod, error) {
	query := `
		SELECT 
			rc.id as category_id,
			rc.name as category_name,
			rc.weight as category_weight
			AVG(r.rating) as avg_score,
			COUNT(r.id) as rating_count,
			DATE(r.created_at, 'weekday 0', '-6 days') as week_start
		FROM ratings r
		JOIN rating_categories rc ON r.rating_category_id = rc.id
		WHERE r.created_at >= ? AND r.created_at <= ?
		GROUP BY rc.id, rc.name, week_start
		ORDER BY rc.name, week_start
	`

	rows, err := r.db.Query(query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query weekly aggregated category scores: %v", err)
	}
	defer rows.Close()

	var scores []models.CategoryRatingOverTimePeriod
	for rows.Next() {
		var score models.CategoryRatingOverTimePeriod
		var weekStart time.Time

		err := rows.Scan(
			&score.CategoryID,
			&score.CategoryName,
			&score.AvgPercent,
			&score.RatingCount,
			&weekStart,
			&score.RatingsTotalCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan weekly category score: %v", err)
		}

		score.Date = weekStart
		scores = append(scores, score)
	}

	return scores, nil
}

func (r *AnalyticsRepository) GetScoresByTicket(startDate, endDate time.Time) ([]models.TicketCategoryScore, error) {
	query := `
		SELECT 
			t.id as ticket_id,
			rc.id as category_id,
			rc.name as category_name,
			AVG(r.rating) as avg_score,
			COUNT(r.id) as rating_count
		FROM ratings r
		JOIN tickets t ON r.ticket_id = t.id
		JOIN rating_categories rc ON r.rating_category_id = rc.id
		WHERE r.created_at >= ? AND r.created_at <= ?
		GROUP BY t.id, rc.id, rc.name
		ORDER BY t.id, rc.name
	`

	rows, err := r.db.Query(query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query scores by ticket: %v", err)
	}
	defer rows.Close()

	var scores []models.TicketCategoryScore
	for rows.Next() {
		var score models.TicketCategoryScore

		err := rows.Scan(
			&score.TicketID,
			&score.CategoryID,
			&score.CategoryName,
			&score.Score,
			&score.RatingCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ticket category score: %v", err)
		}

		scores = append(scores, score)
	}

	return scores, nil
}

func (r *AnalyticsRepository) GetOverallQualityScore(startDate, endDate time.Time) (float64, int, error) {
	query := `
		SELECT 
			AVG(r.rating) as overall_score,
			COUNT(r.id) as total_ratings
		FROM ratings r
		WHERE r.created_at >= ? AND r.created_at <= ?
	`

	var overallScore float64
	var totalRatings int

	err := r.db.QueryRow(query, startDate, endDate).Scan(&overallScore, &totalRatings)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to query overall quality score: %v", err)
	}

	return overallScore, totalRatings, nil
}

func (r *AnalyticsRepository) GetRatingCategories() ([]models.RatingCategory, error) {
	query := `SELECT id, name, weight FROM rating_categories ORDER BY name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query rating categories: %v", err)
	}
	defer rows.Close()

	var categories []models.RatingCategory
	for rows.Next() {
		var category models.RatingCategory

		err := rows.Scan(&category.ID, &category.Name, &category.Weight)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rating category: %v", err)
		}

		categories = append(categories, category)
	}

	return categories, nil
}
