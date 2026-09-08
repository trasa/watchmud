package gameserver

import (
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/player"
)

type HandlerParameter struct {
	Client  Conn
	Player  *player.Player
	Message *message.GameMessage
}

func NewHandlerParameter(c Conn, msg *message.GameMessage) *HandlerParameter {
	return &HandlerParameter{
		Client:  c,
		Player:  c.Player(),
		Message: msg,
	}
}
