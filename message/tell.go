package message

type TellRequest struct {
	Request
	ReceiverPlayerName string
	Value              string
}

type TellResponse struct {
	Response
}

type TellNotification struct {
	Response
	Sender string
	Value  string
}
