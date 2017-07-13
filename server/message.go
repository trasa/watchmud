package server

import "github.com/trasa/watchmud/world"

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
	Sender      string `json:"sender"`
	MessageType string `json:"msg_type"`
	Successful  bool   `json:"success"`
	ResultCode  string `json:"result_code"`
}

type LoginResponse struct {
	MessageType string        `json:"msg_type"`
	Successful  bool          `json:"success"`
	ResultCode  string        `json:"result_code"`
	Player      *world.Player `json:"player"`
}
