package services

import (
	"context"
	"time"

	"github.com/andriawan24/link-short/internal/database"
	"github.com/andriawan24/link-short/internal/utils"
	"github.com/google/uuid"
)

type linkService struct {
	queries *database.Queries
}

type LinkService interface {
	GetTotalCounts(ctx context.Context, userId uuid.UUID, from time.Time, to time.Time) (int64, error)
	GetTotalActiveLinks(ctx context.Context, userId uuid.UUID) (int64, error)
	GetLinks(ctx context.Context, userId uuid.UUID, limit int32, offset int32, orderBy utils.LinkOrderBy) ([]database.GetLinksRow, error)
	GetLink(ctx context.Context, userId uuid.UUID, id uuid.UUID) (database.GetLinkRow, error)
	GetRedirectedLink(ctx context.Context, shortCode string) (string, error)
	InsertLink(ctx context.Context, param database.InsertLinkParams) (database.Link, error)
	DeleteLink(ctx context.Context, param database.DeleteLinkParams) error
}

func NewLinkService(queries *database.Queries) LinkService {
	return &linkService{
		queries: queries,
	}
}

func (l *linkService) GetLink(ctx context.Context, userId uuid.UUID, id uuid.UUID) (database.GetLinkRow, error) {
	param := database.GetLinkParams{
		UserID: userId,
		ID:     id,
	}

	link, err := l.queries.GetLink(ctx, param)
	if err != nil {
		return link, err
	}

	return link, nil
}

func (l *linkService) GetLinks(ctx context.Context, userId uuid.UUID, limit int32, offset int32, orderBy utils.LinkOrderBy) ([]database.GetLinksRow, error) {
	param := database.GetLinksParams{
		UserID:  userId,
		Limit:   limit,
		Offset:  offset,
		OrderBy: orderBy.GetString(),
	}

	links, err := l.queries.GetLinks(ctx, param)
	if err != nil {
		return links, err
	}

	return links, nil
}

func (l *linkService) GetRedirectedLink(ctx context.Context, shortCode string) (string, error) {
	link, err := l.queries.GetRedirectLink(ctx, shortCode)
	if err != nil {
		return link, err
	}

	return link, nil
}

func (l *linkService) InsertLink(ctx context.Context, param database.InsertLinkParams) (database.Link, error) {
	link, err := l.queries.InsertLink(ctx, param)
	if err != nil {
		return link, err
	}

	return link, nil
}

func (l *linkService) GetTotalCounts(ctx context.Context, userId uuid.UUID, from time.Time, to time.Time) (int64, error) {
	param := database.GetTotalClicksParams{
		UserID:   userId,
		FromDate: from,
		ToDate:   to,
	}

	count, err := l.queries.GetTotalClicks(ctx, param)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (l *linkService) GetTotalActiveLinks(ctx context.Context, userId uuid.UUID) (int64, error) {
	count, err := l.queries.GetTotalActiveLinks(ctx, userId)
	if err != nil {
		return count, err
	}

	return count, nil
}

func (l *linkService) DeleteLink(ctx context.Context, param database.DeleteLinkParams) error {
	err := l.queries.DeleteLink(ctx, param)
	return err
}
