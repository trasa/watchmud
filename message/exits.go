package message

type ExitsRequest struct {
	Request
}

type ExitsResponse struct {
	Response
	ExitInfo map[string]string
}
