package world

import (
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/message"
	"log"
)

func (w *World) handlePing(msg *gameserver.HandlerParameter) {
	//noinspection ALL
	if VERBOSE_LOGGING {
		log.Printf("Player %s Ping\n", msg.Player.GetName())
	}
	msg.Player.Send(message.Pong{
		Success:    true,
		ResultCode: "OK",
		Target:     msg.Message.GetPing().Target,
	})
}
