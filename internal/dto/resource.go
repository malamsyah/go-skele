package dto

// SuccessResponse represents success response.
type SuccessResponse struct {
	Data any `json:"data"`
}

// ErrorResponse represents error response.
type ErrorResponse struct {
	Message string `json:"message"`
}

// Resource represents resource data.
type Resource struct {
	ID           uint   `json:"id"`
	Payload      string `json:"payload"`
	ResourceType string `json:"resource_type"`
}
