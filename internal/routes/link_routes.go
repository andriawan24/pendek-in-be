package routes

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/andriawan24/link-short/internal/database"
	"github.com/andriawan24/link-short/internal/models/requests"
	"github.com/andriawan24/link-short/internal/models/responses"
	"github.com/andriawan24/link-short/internal/services"
	"github.com/andriawan24/link-short/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/medama-io/go-useragent"
)

type linkRoutes struct {
	linkService     services.LinkService
	clickLogService services.ClickLogService
	cacheService    services.CacheService
}

func NewLinkRoutes(linkService services.LinkService, clickLogService services.ClickLogService, cacheService services.CacheService) linkRoutes {
	return linkRoutes{
		linkService:     linkService,
		clickLogService: clickLogService,
		cacheService:    cacheService,
	}
}

// GetLink godoc
// @Summary      Get link by ID
// @Description  Get a specific link by its ID
// @Tags         Links
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Link ID"
// @Success      200  {object}  responses.BaseResponse{data=responses.LinkResponse}
// @Failure      400  {object}  responses.ErrorResponse
// @Failure      401  {object}  responses.ErrorResponse
// @Failure      404  {object}  responses.ErrorResponse
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /links/{id} [get]
func (r *linkRoutes) GetLink(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)

	linkId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	link, err := r.linkService.GetLink(ctx.Request.Context(), userId, linkId)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	to := time.Now()
	from := time.Time{}

	totalClicks, err := r.linkService.GetTotalCounts(ctx.Request.Context(), userId, from, to)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	deviceBreakdown, err := r.clickLogService.GetDeviceBreakdownSingleLink(ctx, userId, link.ID, from, to)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	countryBreakdown, err := r.clickLogService.GetTopCountriesSingleLink(ctx, userId, link.ID, from, to)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	devices := responses.MapDeviceBreakdownSingle(deviceBreakdown)
	countries := responses.MapTopCountriesSingle(countryBreakdown)

	utils.RespondOK(ctx, "successfully get link", responses.MapLinkResponse(link, totalClicks, devices, countries))
}

// GetLinks godoc
// @Summary      Get all links
// @Description  Get all links for the authenticated user with pagination
// @Tags         Links
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page     query     int     false  "Page number"       default(1)
// @Param        limit    query     int     false  "Items per page"    default(10)
// @Param        orderBy  query     string  false  "Order by field"    Enums(created_at, counts)
// @Success      200  {object}  responses.BaseResponse{data=[]responses.LinkResponse}
// @Failure      400  {object}  responses.ErrorResponse
// @Failure      401  {object}  responses.ErrorResponse
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /links/all [get]
func (r *linkRoutes) GetLinks(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)
	var (
		page    = 1
		limit   = 10
		orderBy = utils.OrderByCreatedDate
		err     error
	)

	if ctx.Query("page") != "" {
		page, err = strconv.Atoi(ctx.Query("page"))
		if err != nil {
			utils.HandleErrorResponse(ctx, err)
			return
		}
	}

	if ctx.Query("limit") != "" {
		limit, err = strconv.Atoi(ctx.Query("limit"))
		if err != nil {
			utils.HandleErrorResponse(ctx, err)
			return
		}
	}

	if ctx.Query("orderBy") != "" {
		orderBy, err = utils.ParseLinkOrderBy(ctx.Query("orderBy"))
		if err != nil {
			utils.HandleErrorResponse(ctx, err)
			return
		}
	}

	offset := (page - 1) * limit

	links, err := r.linkService.GetLinks(ctx.Request.Context(), userId, int32(limit), int32(offset), orderBy)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.RespondOK(ctx, "successfully get links", responses.MapLinkResponses(links))
}

// InsertLink godoc
// @Summary      Create new link
// @Description  Create a new shortened link
// @Tags         Links
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body requests.InsertLinkParam true "Link details"
// @Success      200  {object}  responses.BaseResponse{data=responses.LinkResponse}
// @Failure      400  {object}  responses.ErrorResponse
// @Failure      401  {object}  responses.ErrorResponse
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /links/create [post]
func (r *linkRoutes) InsertLink(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)

	var body requests.InsertLinkParam

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	param := database.InsertLinkParams{
		OriginalUrl: body.OriginalURL,
		ShortCode:   utils.GenerateShortCode(),
		CustomShortCode: sql.NullString{
			Valid:  body.CustomShortCode != nil,
			String: utils.GetOrElse(body.CustomShortCode, ""),
		},
		UserID: userId,
		ExpiredAt: sql.NullTime{
			Valid: body.ExpiredAt != nil,
			Time:  utils.GetOrElse(body.ExpiredAt, time.Now()),
		},
	}

	link, err := r.linkService.InsertLink(ctx.Request.Context(), param)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.ResponsdJson(ctx, http.StatusCreated, "successfully insert new link", responses.MapLinkDetailResponse(link))
}

// DeleteLink godoc
// @Summary      Delete an existing link
// @Description  Delete a shortened link by its ID
// @Tags         Links
// @Security     BearerAuth
// @Param        id   path      string  true  "Link ID (UUID)"
// @Success      204  {string}  string  "No Content"
// @Failure      400  {object}  responses.ErrorResponse
// @Failure      401  {object}  responses.ErrorResponse
// @Failure      404  {object}  responses.ErrorResponse
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /links/{id} [delete]
func (r *linkRoutes) DeleteLink(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(uuid.UUID)

	linkId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	link, err := r.linkService.GetLink(ctx.Request.Context(), userId, linkId)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	param := database.DeleteLinkParams{
		UserID: userId,
		ID:     link.ID,
	}

	err = r.linkService.DeleteLink(ctx.Request.Context(), param)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	utils.ResponsdJson(ctx, http.StatusNoContent, "successfully insert new link", nil)
}

// Redirect godoc
// @Summary      Redirect to original URL
// @Description  Redirect to the original URL using the short code
// @Tags         Redirect
// @Param        code   path      string  true  "Short code"
// @Success      301  {string}  string  "Redirect to original URL"
// @Failure      404  {object}  responses.ErrorResponse
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /{code} [get]
func (r *linkRoutes) Redirect(ctx *gin.Context) {
	code := ctx.Param("code")
	reqCtx := ctx.Request.Context()

	parser := useragent.NewParser()
	ua := parser.Parse(ctx.Request.UserAgent())

	deviceType := utils.ParseDeviceType(ua)
	country := utils.ParseCountryFromIp(ctx.ClientIP())
	traffic := utils.ParseTrafficSource(ctx.Request.Referer())
	browser := utils.ParseBrowser(ua)

	param := database.InsertClickLogParams{
		Code: code,
		IpAddress: sql.NullString{
			Valid:  ctx.ClientIP() != "",
			String: ctx.ClientIP(),
		},
		UserAgent: sql.NullString{
			Valid:  ctx.Request.UserAgent() != "",
			String: ctx.Request.UserAgent(),
		},
		Referrer: sql.NullString{
			Valid:  ctx.Request.Referer() != "",
			String: ctx.Request.Referer(),
		},
		DeviceType: sql.NullString{
			Valid:  deviceType != "",
			String: deviceType,
		},
		Country: sql.NullString{
			Valid:  country != "",
			String: country,
		},
		Traffic: sql.NullString{
			Valid:  traffic != "",
			String: traffic,
		},
		Browser: sql.NullString{
			Valid:  browser != "",
			String: browser,
		},
	}

	// Try redis
	originalURL, err := r.cacheService.GetURL(reqCtx, code)
	if err == nil && originalURL != "" {
		if _, err := r.clickLogService.InsertClickLog(reqCtx, param); err != nil {
			log.Printf("failed to insert click log (redis hit) for code %s: %v", code, err)
		}
		ctx.Redirect(http.StatusMovedPermanently, originalURL)
		return
	}

	originalURL, err = r.linkService.GetRedirectedLink(reqCtx, code)
	if err != nil {
		utils.HandleErrorResponse(ctx, err)
		return
	}

	go func() {
		_ = r.cacheService.SetURL(context.Background(), code, originalURL, 24*time.Hour)
	}()

	if _, err := r.clickLogService.InsertClickLog(reqCtx, param); err != nil {
		log.Printf("failed to insert click log (db hit) for code %s: %v", code, err)
	}

	ctx.Redirect(http.StatusMovedPermanently, originalURL)
}
