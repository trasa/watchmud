package gameserver

import "github.com/trasa/watchmud/message"

type Instance interface {
	Receive(message *message.IncomingMessage)
}
