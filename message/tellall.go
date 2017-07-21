package message

type TellAllRequest struct {
	Request
	Value string `json:"value"`
}

type TellAllResponse struct {
	Response
}

type TellAllNotification struct {
	Notification
	Value  string `json:"value"`
	Sender string `json:"sender"`
}
