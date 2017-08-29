package message

type TellAllRequest struct {
	Request
	Value string
}

type TellAllResponse struct {
	Response
}

type TellAllNotification struct {
	Notification
	Value  string `json:"value"`
	Sender string `json:"sender"`
}
