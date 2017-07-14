package server


type IncomingMessage struct {
	Client *Client
	Body   map[string]string
}

func newIncomingMessage(client *Client, body map[string]string) *IncomingMessage {
	return &IncomingMessage{
		Client: client,
		Body:   body,
	}
}

type Response struct {
	MessageType string `json:"msg_type"`
	Successful  bool   `json:"success"`
	ResultCode  string `json:"result_code"`
}

type TellAllResponse struct {
	Response
	Sender string `json:"sender"`
}

type LoginResponse struct {
	Response
	Player *Player `json:"player"`
}
