package message

import (
	"github.com/trasa/watchmud/client"
	"github.com/trasa/watchmud/player"
)

type IncomingMessage struct {
	Client  client.Client
	Player  player.Player
	Request Request
}

func New(client client.Client, request Request) *IncomingMessage {
	return &IncomingMessage{
		Client:  client,
		Player:  client.GetPlayer(),
		Request: request,
	}
}
