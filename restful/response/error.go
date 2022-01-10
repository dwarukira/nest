package response

// swagger:model ErrorResponse
type ErrorResponse struct {
	Type       string `json:"type"`
	Title      string `json:"title"`
	Message    string `json:"message"`
	StackTrace string `json:"stack_trace"`
}
