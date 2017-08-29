package message

type TellRequest struct {
	Request
	ReceiverPlayerName string
	Value              string
}

type TellNotification struct {
	Response
	Sender string `json:"sender"`
	Value  string `json:"value"`
}
