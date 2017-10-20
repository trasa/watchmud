package message

type ExitsRequest struct {
	Request
}

type ExitsResponse struct {
	Response
	// maps directionAbbreviation to name of room
	ExitInfo map[string]string
}
