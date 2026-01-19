package services

import (
	"context"

	"github.com/andriawan24/link-short/internal/database"
	"github.com/andriawan24/link-short/internal/models/responses"
)

type dashboardService struct {
	queries *database.Queries
}

type DashboardService interface {
	GetLandingStats(ctx context.Context) (responses.LandingStatsResponse, error)
}

func NewDashboardService(queries *database.Queries) DashboardService {
	return &dashboardService{
		queries: queries,
	}
}

func (s *dashboardService) GetLandingStats(ctx context.Context) (responses.LandingStatsResponse, error) {
	var stats responses.LandingStatsResponse

	totalLinks, err := s.queries.GetTotalLinksCreated(ctx)
	if err != nil {
		return stats, err
	}

	totalUsers, err := s.queries.GetTotalActiveUsers(ctx)
	if err != nil {
		return stats, err
	}

	totalClicks, err := s.queries.GetGlobalTotalClicks(ctx)
	if err != nil {
		return stats, err
	}

	stats.TotalLinks = totalLinks
	stats.TotalActiveUsers = totalUsers
	stats.TotalClicks = totalClicks

	return stats, nil
}
