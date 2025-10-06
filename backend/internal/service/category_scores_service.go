package service

import (
	"sort"
	"time"

	"go-grpc-backend/internal/models"
	"go-grpc-backend/internal/repository"
	"go-grpc-backend/proto"

	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// GetAggregatedCategoryScores retrieves and aggregates category scores over time
// It automatically selects daily or weekly granularity based on the date range
func GetAggregatedCategoryScores(repo repository.AnalyticsRepositoryInterface, startDate, endDate time.Time) (*proto.AggregatedCategoryScoresResponse, error) {
	duration := endDate.Sub(startDate)
	useWeekly := duration > 30*24*time.Hour

	var (
		rows []models.CategoryRatingOverTimePeriod
		err  error
	)
	if useWeekly {
		rows, err = repo.GetWeeklyAggregatedCategoryRatings(startDate, endDate)
	} else {
		rows, err = repo.GetDailyAggregatedCategoryRatings(startDate, endDate)
	}
	if err != nil {
		return nil, err
	}

	// Group by category → collect series slice
	byCat := make(map[int32]*proto.CategorySeries)
	for _, r := range rows {
		cid := int32(r.CategoryID)

		series, ok := byCat[cid]
		if !ok {
			series = &proto.CategorySeries{
				CategoryId:         cid,
				CategoryName:       r.CategoryName,
				CategoryTotalCount: 0,
				Scores:             nil,
			}
			byCat[cid] = series
		}

		score := CalculateCategoryScore(r.AvgPercent, r.CategoryWeight)
		series.Scores = append(series.Scores, &proto.ScorePoint{
			Date:  timestamppb.New(r.Date),
			Score: float32(score),
			Count: wrapperspb.Int32(int32(r.RatingCount)),
		})
		series.CategoryTotalCount += int32(r.RatingCount)
	}

	categories := make([]*proto.CategorySeries, 0, len(byCat))
	for _, s := range byCat {
		sort.Slice(s.Scores, func(i, j int) bool {
			return s.Scores[i].Date.AsTime().Before(s.Scores[j].Date.AsTime())
		})
		categories = append(categories, s)
	}

	gran := proto.Granularity_GRANULARITY_DAY
	if useWeekly {
		gran = proto.Granularity_GRANULARITY_WEEK
	}

	resp := &proto.AggregatedCategoryScoresResponse{
		Granularity: gran,
		BucketRange: &proto.BucketRange{
			Start: timestamppb.New(startDate),
			End:   timestamppb.New(endDate),
		},
		Categories: categories,
	}
	return resp, nil
}
