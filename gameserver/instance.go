package gameserver

import (
	"github.com/trasa/watchmud/client"
)

type Instance interface {
	Receive(handlerParam *HandlerParameter)
	Logout(c client.Client, cause string)
}
