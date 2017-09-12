package message

import "fmt"

type SayRequest struct {
	Request
	Value string `json:"value"`
}

type SayResponse struct {
	Response
	Value string `json:"value"`
}

func (r SayResponse) String() string {
	return fmt.Sprintf("[Say msgType=%s, result_code=%s, value=%s]", r.GetMessageType(), r.GetResultCode(), r.Value)
}

type SayNotification struct {
	Response
	Sender string `json:"sender"`
	Value  string `json:"value"`
}
