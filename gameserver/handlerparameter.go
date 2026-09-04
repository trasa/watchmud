package gameserver

import (
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/client"
	"github.com/trasa/watchmud/player"
)

type HandlerParameter struct {
	Client  client.Client
	Player  *player.Player
	Message *message.GameMessage
}

func NewHandlerParameter(c client.Client, msg *message.GameMessage) *HandlerParameter {
	return &HandlerParameter{
		Client:  c,
		Player:  c.Player(),
		Message: msg,
	}
}
