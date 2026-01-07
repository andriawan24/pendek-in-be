package requests

import "time"

type InsertLinkParam struct {
	OriginalURL     string     `json:"original_url" binding:"required"`
	CustomShortCode *string    `json:"custom_short_code"`
	ExpiredAt       *time.Time `json:"expired_at"`
}
