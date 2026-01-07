package requests

type RefreshParam struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
