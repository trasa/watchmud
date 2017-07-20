package message

import (
	"github.com/trasa/watchmud/client"
	"github.com/trasa/watchmud/player"
)

type IncomingMessage struct {
	Client client.Client
	Player player.Player
	Body   map[string]string
}

func NewIncomingMessage(client client.Client, body map[string]string) *IncomingMessage {
	return &IncomingMessage{
		Client: client,
		Player: client.GetPlayer(),
		Body:   body,
	}
}
