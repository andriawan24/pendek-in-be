package responses

import (
	"time"

	"github.com/andriawan24/link-short/internal/database"
)

type DashboardResponse struct {
	TotalClicks      int64              `json:"total_clicks"`
	TotalActiveLinks int64              `json:"total_active_links"`
	TopLink          *TopLink           `json:"top_link"`
	Overviews        []AnalyticOverview `json:"overviews"`
	Recents          []LinkResponse     `json:"recents"`
}

type AnalyticOverview struct {
	Date  time.Time `json:"date"`
	Value int       `json:"value"`
}

type TopLink struct {
	Link        LinkResponse `json:"link"`
	TotalClicks int64        `json:"total_clicks"`
}

func MapAnalyticsResponse(rows []database.GetByDateRangeRow) []AnalyticOverview {
	var overviews []AnalyticOverview

	for _, item := range rows {
		overviews = append(overviews, AnalyticOverview{
			Date:  item.Date,
			Value: int(item.TotalClick),
		})
	}

	return overviews
}
