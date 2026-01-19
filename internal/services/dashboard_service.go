package services

import (
	"context"

	"github.com/andriawan24/link-short/internal/database"
	"github.com/andriawan24/link-short/internal/models/responses"
)

type dashboardService struct {
	ctx     context.Context
	queries *database.Queries
}

type DashboardService interface {
	GetLandingStats() (responses.LandingStatsResponse, error)
}

func NewDashboardService(ctx context.Context, queries *database.Queries) DashboardService {
	return &dashboardService{
		ctx:     ctx,
		queries: queries,
	}
}

func (s *dashboardService) GetLandingStats() (responses.LandingStatsResponse, error) {
	var stats responses.LandingStatsResponse

	totalLinks, err := s.queries.GetTotalLinksCreated(s.ctx)
	if err != nil {
		return stats, err
	}

	totalUsers, err := s.queries.GetTotalActiveUsers(s.ctx)
	if err != nil {
		return stats, err
	}

	totalClicks, err := s.queries.GetGlobalTotalClicks(s.ctx)
	if err != nil {
		return stats, err
	}

	stats.TotalLinks = totalLinks
	stats.TotalActiveUsers = totalUsers
	stats.TotalClicks = totalClicks

	return stats, nil
}
