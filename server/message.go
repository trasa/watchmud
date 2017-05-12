package server

type Message struct {
	Client *Client
	Body   map[string]string
}

func newMessage(client *Client, body map[string]string) *Message {
	return &Message{
		Client: client,
		Body:   body,
	}
}
