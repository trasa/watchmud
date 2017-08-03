package message

// We received a request, but it was incorrectly formatted or otherwise problematic.
type ErrorRequest struct {
	Request
	Error error
}

type ErrorResponse struct {
	Response
	Error string `json:"error"`
}
