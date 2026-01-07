package responses

type BaseResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Error   any    `json:"error,omitempty"`
}
