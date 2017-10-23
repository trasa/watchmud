package world

import (
	"github.com/trasa/watchmud/direction"
	"github.com/trasa/watchmud/message"
)

func (w *World) handleExits(msg *message.IncomingMessage) {
	r := w.getRoomContainingPlayer(msg.Player)
	if r == nil {
		r = w.voidRoom
	}
	// convert directions to strings because json
	messageExitInfo := make(map[string]string)
	for k, v := range r.GetExitInfo(false) {
		s, _ := direction.DirectionToAbbreviation(k)
		messageExitInfo[s] = v
	}
	resp := message.ExitsResponse{
		Response: message.NewSuccessfulResponse("exits"),
		ExitInfo: messageExitInfo,
	}
	msg.Player.Send(resp)
}
