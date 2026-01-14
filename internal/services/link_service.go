package services

import (
	"context"
	"time"

	"github.com/andriawan24/link-short/internal/database"
	"github.com/andriawan24/link-short/internal/utils"
	"github.com/google/uuid"
)

type linkService struct {
	ctx     context.Context
	queries *database.Queries
}

type LinkService interface {
	GetTotalCounts(userId uuid.UUID, from time.Time, to time.Time) (int64, error)
	GetTotalActiveLinks(userId uuid.UUID) (int64, error)
	GetLinks(userId uuid.UUID, limit int32, offset int32, orderBy utils.LinkOrderBy) ([]database.GetLinksRow, error)
	GetLink(userId uuid.UUID, id uuid.UUID) (database.Link, error)
	GetRedirectedLink(shortCode string) (string, error)
	InsertLink(param database.InsertLinkParams) (database.Link, error)
	DeleteLink(param database.DeleteLinkParams) error
}

func NewLinkService(ctx context.Context, queries *database.Queries) LinkService {
	return &linkService{
		ctx:     ctx,
		queries: queries,
	}
}

func (l *linkService) GetLink(userId uuid.UUID, id uuid.UUID) (database.Link, error) {
	param := database.GetLinkParams{
		UserID: userId,
		ID:     id,
	}

	link, err := l.queries.GetLink(l.ctx, param)
	if err != nil {
		return link, err
	}

	return link, nil
}

func (l *linkService) GetLinks(userId uuid.UUID, limit int32, offset int32, orderBy utils.LinkOrderBy) ([]database.GetLinksRow, error) {
	param := database.GetLinksParams{
		UserID:  userId,
		Limit:   limit,
		Offset:  offset,
		OrderBy: orderBy.GetString(),
	}

	links, err := l.queries.GetLinks(l.ctx, param)
	if err != nil {
		return links, err
	}

	return links, nil
}

func (l *linkService) GetRedirectedLink(shortCode string) (string, error) {
	link, err := l.queries.GetRedirectLink(l.ctx, shortCode)
	if err != nil {
		return link, err
	}

	return link, nil
}

func (l *linkService) InsertLink(param database.InsertLinkParams) (database.Link, error) {
	link, err := l.queries.InsertLink(l.ctx, param)
	if err != nil {
		return link, err
	}

	return link, nil
}

func (l *linkService) GetTotalCounts(userId uuid.UUID, from time.Time, to time.Time) (int64, error) {
	param := database.GetTotalClicksParams{
		UserID:   userId,
		FromDate: from,
		ToDate:   to,
	}

	count, err := l.queries.GetTotalClicks(l.ctx, param)
	if err != nil {
		return 0, err
	}

	total, ok := count.(int64)
	if !ok {
		return 0, nil
	}

	return total, nil
}

func (l *linkService) GetTotalActiveLinks(userId uuid.UUID) (int64, error) {
	count, err := l.queries.GetTotalActiveLinks(l.ctx, userId)
	if err != nil {
		return count, err
	}

	return count, nil
}

func (l *linkService) DeleteLink(param database.DeleteLinkParams) error {
	err := l.queries.DeleteLink(l.ctx, param)
	return err
}
