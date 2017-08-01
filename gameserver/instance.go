package gameserver

import (
	"github.com/trasa/watchmud/client"
	"github.com/trasa/watchmud/message"
)

type Instance interface {
	Receive(message *message.IncomingMessage)
	Logout(c client.Client, cause string)
}
