package routes

import (
	"github.com/andriawan24/link-short/internal/services"
	"github.com/andriawan24/link-short/internal/utils"
	"github.com/gin-gonic/gin"
)

type dashboardRoutes struct {
	dashboardService services.DashboardService
}

func NewDashboardRoutes(dashboardService services.DashboardService) dashboardRoutes {
	return dashboardRoutes{
		dashboardService: dashboardService,
	}
}

// GetLandingStats godoc
// @Summary      Get landing page statistics
// @Description  Get total links created, active users, and total clicks
// @Tags         Dashboard
// @Accept       json
// @Produce      json
// @Success      200  {object}  responses.BaseResponse{data=responses.LandingStatsResponse}
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /dashboard/stats [get]
func (r *dashboardRoutes) GetLandingStats(ctx *gin.Context) {
	stats, err := r.dashboardService.GetLandingStats(ctx.Request.Context())
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.RespondOK(ctx, "successfully get landing stats", stats)
}
