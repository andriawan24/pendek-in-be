package routes

import (
	"time"

	"github.com/andriawan24/link-short/internal/models/responses"
	"github.com/andriawan24/link-short/internal/services"
	"github.com/andriawan24/link-short/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type analyticRoutes struct {
	linkService     services.LinkService
	clickLogService services.ClickLogService
}

func NewAnalyticRoutes(linkService services.LinkService, clickLogService services.ClickLogService) analyticRoutes {
	return analyticRoutes{
		linkService:     linkService,
		clickLogService: clickLogService,
	}
}

func (r *analyticRoutes) GetDashboard(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)

	from := time.Time{}
	to := time.Now()

	totalClicks, err := r.linkService.GetTotalCounts(userId, from, to)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	totalActiveLinks, err := r.linkService.GetTotalActiveLinks(userId)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	topLinks, err := r.linkService.GetLinks(userId, 1, 0, utils.OrderByCounts)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	overviews, err := r.clickLogService.GetByDateRange(ctx, userId, from, to)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	recents, err := r.linkService.GetLinks(userId, 5, 0, utils.OrderByCreatedDate)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}
	recentResponse := responses.MapLinkResponses(recents)

	response := responses.DashboardResponse{
		TotalClicks:      totalClicks,
		TotalActiveLinks: totalActiveLinks,
		Overviews:        responses.MapAnalyticsResponse(overviews),
		Recents:          recentResponse,
	}

	if len(topLinks) > 0 {
		linkResponse := responses.MapLinkResponses(topLinks)
		response.TopLink = &responses.TopLink{
			Link:        linkResponse[0],
			TotalClicks: linkResponse[0].ClickCount,
		}
	}

	utils.RespondOK(ctx, "successfully get dashboard", response)
}

func (r *analyticRoutes) GetAnalytics(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)

	rangeParam := ctx.DefaultQuery("range", "30d")
	timeRange := utils.ParseTimeRange(rangeParam)

	to := time.Now()
	from := timeRange.GetFromDate()

	overviews, err := r.clickLogService.GetByDateRange(ctx, userId, from, to)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	totalClicks, err := r.linkService.GetTotalCounts(userId, from, to)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	totalActiveLinks, err := r.linkService.GetTotalActiveLinks(userId)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	topLinks, err := r.linkService.GetLinks(userId, 1, 0, utils.OrderByCounts)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	deviceBreakdown, err := r.clickLogService.GetDeviceBreakdown(ctx, userId, from, to)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	topCountries, err := r.clickLogService.GetTopCountries(ctx, userId, from, to)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	trafficSources, err := r.clickLogService.GetTrafficSources(ctx, userId, from, to)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	browserUsage, err := r.clickLogService.GetBrowserUsage(ctx, userId, from, to)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	// Calculate average daily clicks
	daysDiff := to.Sub(from).Hours() / 24
	if daysDiff < 1 {
		daysDiff = 1
	}
	avgDailyClick := totalClicks / int64(daysDiff)

	response := responses.AnalyticsResponse{
		TimeRange:        string(timeRange),
		FromDate:         from,
		ToDate:           to,
		TotalClicks:      totalClicks,
		TotalActiveLinks: totalActiveLinks,
		AvgDailyClick:    avgDailyClick,
		Overviews:        responses.MapAnalyticsResponse(overviews),
		DeviceBreakdowns: responses.MapDeviceBreakdown(deviceBreakdown),
		TopCountries:     responses.MapTopCountries(topCountries),
		TrafficSources:   responses.MapTrafficSources(trafficSources),
		BrowserUsages:    responses.MapBrowserUsage(browserUsage),
	}

	if len(topLinks) > 0 {
		linkResponse := responses.MapLinkResponses(topLinks)
		response.TopLink = &responses.TopLink{
			Link:        linkResponse[0],
			TotalClicks: linkResponse[0].ClickCount,
		}
	}

	utils.RespondOK(ctx, "successfully get analytics", response)
}
