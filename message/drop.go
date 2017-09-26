package message

type DropRequest struct {
	Request
	Target string
}

type DropResponse struct {
	Response
}
