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

	totalClicks, err := r.linkService.GetTotalCounts(userId)
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

	from := time.Now().AddDate(0, -1, 0)
	to := time.Now()
	overviews, err := r.clickLogService.GetByDateRange(ctx, userId, from, to)

	recents, err := r.linkService.GetLinks(userId, 5, 0, utils.OrderByCreatedDate)
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
