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

type AnalyticsResponse struct {
	TimeRange        string             `json:"time_range"`
	FromDate         time.Time          `json:"from_date"`
	ToDate           time.Time          `json:"to_date"`
	TotalClicks      int64              `json:"total_clicks"`
	TotalActiveLinks int64              `json:"total_active_links"`
	TopLink          *TopLink           `json:"top_link"`
	AvgDailyClick    int64              `json:"avg_daily_click"`
	Overviews        []AnalyticOverview `json:"overviews"`
	DeviceBreakdowns []TypeValue        `json:"device_breakdowns"`
	TopCountries     []TypeValue        `json:"top_countries"`
	TrafficSources   []TypeValue        `json:"traffic_sources"`
	BrowserUsages    []TypeValue        `json:"browser_usages"`
}

type TypeValue struct {
	Type  string `json:"type"`
	Value int64  `json:"value"`
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

func MapDeviceBreakdown(rows []database.GetDeviceBreakdownRow) []TypeValue {
	var result []TypeValue

	for _, item := range rows {
		result = append(result, TypeValue{
			Type:  item.DeviceType,
			Value: item.Total,
		})
	}

	return result
}

func MapDeviceBreakdownSingle(rows []database.GetDeviceBreakdownSingleRow) []TypeValue {
	var result []TypeValue

	for _, item := range rows {
		result = append(result, TypeValue{
			Type:  item.DeviceType,
			Value: item.Total,
		})
	}

	return result
}

func MapTopCountries(rows []database.GetTopCountriesRow) []TypeValue {
	var result []TypeValue

	for _, item := range rows {
		result = append(result, TypeValue{
			Type:  item.Country,
			Value: item.Total,
		})
	}

	return result
}

func MapTopCountriesSingle(rows []database.GetTopCountriesSingleRow) []TypeValue {
	var result []TypeValue

	for _, item := range rows {
		result = append(result, TypeValue{
			Type:  item.Country,
			Value: item.Total,
		})
	}

	return result
}

func MapTrafficSources(rows []database.GetTrafficSourcesRow) []TypeValue {
	var result []TypeValue

	for _, item := range rows {
		result = append(result, TypeValue{
			Type:  item.TrafficSource,
			Value: item.Total,
		})
	}

	return result
}

func MapBrowserUsage(rows []database.GetBrowserUsageRow) []TypeValue {
	var result []TypeValue

	for _, item := range rows {
		result = append(result, TypeValue{
			Type:  item.Browser,
			Value: item.Total,
		})
	}

	return result
}
