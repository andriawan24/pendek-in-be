package services

import (
	"context"
	"time"

	"github.com/andriawan24/link-short/internal/database"
	"github.com/google/uuid"
)

type clickLogService struct {
	queries *database.Queries
}

type ClickLogService interface {
	InsertClickLog(ctx context.Context, param database.InsertClickLogParams) (database.ClickLog, error)
	GetByDateRange(ctx context.Context, userId uuid.UUID, from time.Time, to time.Time) ([]database.GetByDateRangeRow, error)
	GetDeviceBreakdown(ctx context.Context, userId uuid.UUID, from time.Time, to time.Time) ([]database.GetDeviceBreakdownRow, error)
	GetDeviceBreakdownSingleLink(ctx context.Context, userId uuid.UUID, linkId uuid.UUID, from time.Time, to time.Time) ([]database.GetDeviceBreakdownSingleRow, error)
	GetTopCountries(ctx context.Context, userId uuid.UUID, from time.Time, to time.Time) ([]database.GetTopCountriesRow, error)
	GetTopCountriesSingleLink(ctx context.Context, userId uuid.UUID, linkId uuid.UUID, from time.Time, to time.Time) ([]database.GetTopCountriesSingleRow, error)
	GetTrafficSources(ctx context.Context, userId uuid.UUID, from time.Time, to time.Time) ([]database.GetTrafficSourcesRow, error)
	GetBrowserUsage(ctx context.Context, userId uuid.UUID, from time.Time, to time.Time) ([]database.GetBrowserUsageRow, error)
}

func NewClickLogService(queries *database.Queries) ClickLogService {
	return &clickLogService{
		queries: queries,
	}
}

func (c *clickLogService) InsertClickLog(ctx context.Context, param database.InsertClickLogParams) (database.ClickLog, error) {
	clickLog, err := c.queries.InsertClickLog(ctx, param)
	if err != nil {
		return clickLog, err
	}

	return clickLog, nil
}

func (c *clickLogService) GetByDateRange(ctx context.Context, userId uuid.UUID, from time.Time, to time.Time) ([]database.GetByDateRangeRow, error) {
	logs, err := c.queries.GetByDateRange(ctx, database.GetByDateRangeParams{
		FromDate: from,
		ToDate:   to,
		UserID:   userId,
	})
	if err != nil {
		return logs, err
	}

	return logs, nil
}

func (c *clickLogService) GetDeviceBreakdown(ctx context.Context, userId uuid.UUID, from time.Time, to time.Time) ([]database.GetDeviceBreakdownRow, error) {
	devices, err := c.queries.GetDeviceBreakdown(ctx, database.GetDeviceBreakdownParams{
		FromDate: from,
		ToDate:   to,
		UserID:   userId,
	})
	if err != nil {
		return nil, err
	}

	return devices, nil
}

func (c *clickLogService) GetTopCountries(ctx context.Context, userId uuid.UUID, from time.Time, to time.Time) ([]database.GetTopCountriesRow, error) {
	countries, err := c.queries.GetTopCountries(ctx, database.GetTopCountriesParams{
		FromDate: from,
		ToDate:   to,
		UserID:   userId,
	})
	if err != nil {
		return nil, err
	}

	return countries, nil
}

func (c *clickLogService) GetTrafficSources(ctx context.Context, userId uuid.UUID, from time.Time, to time.Time) ([]database.GetTrafficSourcesRow, error) {
	sources, err := c.queries.GetTrafficSources(ctx, database.GetTrafficSourcesParams{
		FromDate: from,
		ToDate:   to,
		UserID:   userId,
	})
	if err != nil {
		return nil, err
	}

	return sources, nil
}

func (c *clickLogService) GetBrowserUsage(ctx context.Context, userId uuid.UUID, from time.Time, to time.Time) ([]database.GetBrowserUsageRow, error) {
	browsers, err := c.queries.GetBrowserUsage(ctx, database.GetBrowserUsageParams{
		FromDate: from,
		ToDate:   to,
		UserID:   userId,
	})
	if err != nil {
		return nil, err
	}

	return browsers, nil
}

func (c *clickLogService) GetDeviceBreakdownSingleLink(ctx context.Context, userId uuid.UUID, linkId uuid.UUID, from time.Time, to time.Time) ([]database.GetDeviceBreakdownSingleRow, error) {
	devices, err := c.queries.GetDeviceBreakdownSingle(ctx, database.GetDeviceBreakdownSingleParams{
		FromDate: from,
		ToDate:   to,
		UserID:   userId,
		ID:       linkId,
	})
	if err != nil {
		return nil, err
	}

	return devices, nil
}

func (c *clickLogService) GetTopCountriesSingleLink(ctx context.Context, userId uuid.UUID, linkId uuid.UUID, from time.Time, to time.Time) ([]database.GetTopCountriesSingleRow, error) {
	countries, err := c.queries.GetTopCountriesSingle(ctx, database.GetTopCountriesSingleParams{
		FromDate: from,
		ToDate:   to,
		UserID:   userId,
		ID:       linkId,
	})
	if err != nil {
		return nil, err
	}

	return countries, nil
}
