package message

type TellRequest struct {
	Request
	ReceiverPlayerName string `json:"receiver"`
	Value              string `json:"value"`
}

type TellResponse struct {
	Response
}

type TellNotification struct {
	Notification
	Sender string `json:"sender"`
	Value  string `json:"value"`
}
