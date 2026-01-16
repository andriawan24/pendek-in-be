package responses

import (
	"time"

	"github.com/andriawan24/link-short/internal/database"
	"github.com/google/uuid"
)

type LinkResponse struct {
	ID               uuid.UUID   `json:"id"`
	OriginalURL      string      `json:"original_url"`
	ShortCode        string      `json:"short_code"`
	CustomShortCode  *string     `json:"custom_short_code"`
	ClickCount       int64       `json:"click_count"`
	ExpiredAt        *time.Time  `json:"expired_at"`
	CreatedAt        time.Time   `json:"created_at"`
	DeviceBreakdowns []TypeValue `json:"device_breakdowns"`
	TopCountries     []TypeValue `json:"top_countries"`
}

func MapLinkResponses(links []database.GetLinksRow) []LinkResponse {
	response := make([]LinkResponse, len(links))

	for idx, link := range links {
		var customShortCode *string = nil
		if link.CustomShortCode.Valid {
			customShortCode = &link.CustomShortCode.String
		}

		var expiredAt *time.Time = nil
		if link.ExpiredAt.Valid {
			expiredAt = &link.ExpiredAt.Time
		}

		response[idx] = LinkResponse{
			ID:              link.ID,
			OriginalURL:     link.OriginalUrl,
			ShortCode:       link.ShortCode,
			CustomShortCode: customShortCode,
			ExpiredAt:       expiredAt,
			ClickCount:      link.Counts,
			CreatedAt:       link.CreatedAt,
		}
	}

	return response
}

func MapLinkResponse(link database.GetLinkRow, totalClicks int64, devices []TypeValue, countries []TypeValue) LinkResponse {
	var customShortCode *string = nil
	if link.CustomShortCode.Valid {
		customShortCode = &link.CustomShortCode.String
	}

	var expiredAt *time.Time = nil
	if link.ExpiredAt.Valid {
		expiredAt = &link.ExpiredAt.Time
	}

	response := LinkResponse{
		ID:               link.ID,
		OriginalURL:      link.OriginalUrl,
		ShortCode:        link.ShortCode,
		CustomShortCode:  customShortCode,
		ExpiredAt:        expiredAt,
		CreatedAt:        link.CreatedAt,
		ClickCount:       totalClicks,
		DeviceBreakdowns: devices,
		TopCountries:     countries,
	}

	return response
}

func MapLinkDetailResponse(link database.Link) LinkResponse {
	var customShortCode *string = nil
	if link.CustomShortCode.Valid {
		customShortCode = &link.CustomShortCode.String
	}

	var expiredAt *time.Time = nil
	if link.ExpiredAt.Valid {
		expiredAt = &link.ExpiredAt.Time
	}

	response := LinkResponse{
		ID:              link.ID,
		OriginalURL:     link.OriginalUrl,
		ShortCode:       link.ShortCode,
		CustomShortCode: customShortCode,
		ExpiredAt:       expiredAt,
		CreatedAt:       link.CreatedAt,
	}

	return response
}
