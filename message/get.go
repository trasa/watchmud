package message

type GetRequest struct {
	Request
	Targets []string
}

type GetResponse struct {
	Response
}

type GetNotification struct {
	Response
	PlayerName string
	Target     string
}
