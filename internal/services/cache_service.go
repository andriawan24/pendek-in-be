package services

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type cacheService struct {
	rdb    *redis.Client
	prefix string
}

type CacheService interface {
	GetURL(ctx context.Context, code string) (string, error)
	SetURL(ctx context.Context, code, originalURL string, ttl time.Duration) error
	InvalidateURL(ctx context.Context, code string) error
}

func NewCacheService(rdb *redis.Client) CacheService {
	return &cacheService{
		rdb:    rdb,
		prefix: "url:",
	}
}

func (c *cacheService) GetURL(ctx context.Context, code string) (string, error) {
	return c.rdb.Get(ctx, c.prefix+code).Result()
}

func (c *cacheService) InvalidateURL(ctx context.Context, code string) error {
	return c.rdb.Del(ctx, c.prefix+code).Err()
}

func (c *cacheService) SetURL(ctx context.Context, code string, originalURL string, ttl time.Duration) error {
	return c.rdb.Set(ctx, c.prefix+code, originalURL, ttl).Err()
}
