package responses

import (
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	IsActive        bool      `json:"is_active"`
	IsVerified      bool      `json:"is_verified"`
	ProfileImageUrl string    `json:"profile_image_url"`
}

type LoginResponse struct {
	Token                 string       `json:"token"`
	TokenExpiredAt        time.Time    `json:"token_expired_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiredAt time.Time    `json:"refresh_token_expired_at"`
	User                  UserResponse `json:"user"`
	AuthURL               string       `json:"auth_url"`
	State                 string       `json:"state"`
}
