package message

type DropRequest struct {
	Request
	Target string
}

type DropResponse struct {
	Response
}

type DropNotification struct {
	Response
	PlayerName string
	Target     string
}
