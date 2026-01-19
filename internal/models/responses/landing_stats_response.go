package responses

type LandingStatsResponse struct {
	TotalLinks       int64 `json:"total_links"`
	TotalActiveUsers int64 `json:"total_active_users"`
	TotalClicks      int64 `json:"total_clicks"`
}
