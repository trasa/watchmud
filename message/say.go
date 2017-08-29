package message

type SayRequest struct {
	Request
	Value string `json:"value"`
}

type SayResponse struct {
	Response
	Value string `json:"value"`
}

type SayNotification struct {
	Response
	Sender string `json:"sender"`
	Value  string `json:"value"`
}
