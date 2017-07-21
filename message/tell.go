package message

type TellRequest struct {
	Request
	ReceiverPlayerName string `json:"receiver"`
	Value              string `json:"value"`
}

type TellResponse struct {
	Response
}
