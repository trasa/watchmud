package world

import (
	"github.com/trasa/watchmud/message"
	"log"
)

func (w *World) handleNoOp(msg *message.IncomingMessage) {
	//noinspection ALL
	if VERBOSE_LOGGING {
		log.Printf("Player %s sent a No Op Request\n", msg.Player.GetName())
	}
	msg.Player.Send(message.NoOpResponse{
		Response: message.NewSuccessfulResponse("no_op"),
	})
}
