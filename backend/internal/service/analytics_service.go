package service

import (
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
