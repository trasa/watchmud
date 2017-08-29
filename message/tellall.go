package message

type TellAllRequest struct {
	Request
	Value string
}

type TellAllResponse struct {
	Response
}

type TellAllNotification struct {
	Response
	Value  string `json:"value"`
	Sender string `json:"sender"`
}
