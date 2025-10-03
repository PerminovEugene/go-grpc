package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"go-grpc-backend/internal/database"
	"go-grpc-backend/internal/repository"
	"go-grpc-backend/internal/service"
	"go-grpc-backend/proto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AnalyticsServer struct {
	proto.UnimplementedAnalyticsServiceServer
	analyticsRepo *repository.AnalyticsRepository
	scoreService  *service.ScoreService
	grpcServer    *grpc.Server
}

func NewAnalyticsServer() (*AnalyticsServer, error) {
	db, err := database.NewDatabase()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %v", err)
	}

	analyticsRepo := repository.NewAnalyticsRepository(db.DB)
	scoreService := service.NewScoreService(analyticsRepo)
	grpcServer := grpc.NewServer()

	server := &AnalyticsServer{
		analyticsRepo: analyticsRepo,
		scoreService:  scoreService,
		grpcServer:    grpcServer,
	}

	proto.RegisterAnalyticsServiceServer(grpcServer, server)

	return server, nil
}

func (s *AnalyticsServer) Start(port string) error {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to create listener: %v", err)
	}

	log.Printf("Starting Analytics gRPC server on port %s", port)

	
	if err := s.grpcServer.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

func (s *AnalyticsServer) Stop() {
	log.Println("Stopping Analytics gRPC server...")
	s.grpcServer.GracefulStop()
}

func (s *AnalyticsServer) GetAggregatedCategoryScores(ctx context.Context, req *proto.AggregatedCategoryScoresRequest) (*proto.DailyAggregatedScoresResponse, error) {
	startDate := req.StartDate.AsTime()
	endDate := req.EndDate.AsTime()
	
	scores, err := s.scoreService.GetAggregatedCategoryScores(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get aggregated category scores: %v", err)
	}
	
	var protoScores []*proto.CategoryScore
	for _, score := range scores {
		protoScore := &proto.CategoryScore{
			CategoryId:   int32(score.CategoryID),
			CategoryName: score.CategoryName,
			Score:        score.Score,
			RatingCount:  int32(score.RatingCount),
			Date:         timestamppb.New(score.Date),
		}
		protoScores = append(protoScores, protoScore)
	}
	
	return &proto.DailyAggregatedScoresResponse{
		Scores:    protoScores,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}, nil
}

func (s *AnalyticsServer) GetScoresByTicket(ctx context.Context, req *proto.ScoresByTicketRequest) (*proto.ScoresByTicketResponse, error) {
	startDate := req.StartDate.AsTime()
	endDate := req.EndDate.AsTime()
	
	scores, err := s.scoreService.GetScoresByTicket(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get scores by ticket: %v", err)
	}
	
	var protoScores []*proto.TicketCategoryScore
	for _, score := range scores {
		protoScore := &proto.TicketCategoryScore{
			TicketId:     int32(score.TicketID),
			CategoryId:   int32(score.CategoryID),
			CategoryName: score.CategoryName,
			Score:        score.Score,
			RatingCount:  int32(score.RatingCount),
		}
		protoScores = append(protoScores, protoScore)
	}
	
	return &proto.ScoresByTicketResponse{
		Scores:    protoScores,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}, nil
}

func (s *AnalyticsServer) GetOverallQualityScore(ctx context.Context, req *proto.OverallQualityScoreRequest) (*proto.OverallQualityScoreResponse, error) {
	startDate := req.StartDate.AsTime()
	endDate := req.EndDate.AsTime()
	
	overallScore, totalRatings, err := s.scoreService.GetOverallQualityScore(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get overall quality score: %v", err)
	}
	
	return &proto.OverallQualityScoreResponse{
		OverallScore: overallScore,
		TotalRatings: int32(totalRatings),
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
	}, nil
}

func (s *AnalyticsServer) GetPeriodOverPeriodChange(ctx context.Context, req *proto.PeriodOverPeriodChangeRequest) (*proto.PeriodOverPeriodChangeResponse, error) {
	currentStart := req.CurrentStart.AsTime()
	currentEnd := req.CurrentEnd.AsTime()
	previousStart := req.PreviousStart.AsTime()
	previousEnd := req.PreviousEnd.AsTime()
	
	currentScore, previousScore, err := s.scoreService.GetPeriodOverPeriodChange(
		currentStart, currentEnd, previousStart, previousEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to get period over period change: %v", err)
	}
	
	changePercentage := s.scoreService.CalculateChangePercentage(currentScore, previousScore)
	
	return &proto.PeriodOverPeriodChangeResponse{
		CurrentPeriodScore:  currentScore,
		PreviousPeriodScore: previousScore,
		ChangePercentage:    changePercentage,
		CurrentStart:        req.CurrentStart,
		CurrentEnd:          req.CurrentEnd,
		PreviousStart:       req.PreviousStart,
		PreviousEnd:         req.PreviousEnd,
	}, nil
}
