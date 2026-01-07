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
