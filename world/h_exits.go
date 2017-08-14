package world

import "github.com/trasa/watchmud/message"

func (w *World) handleExits(msg *message.IncomingMessage) {
	r := w.getRoomContainingPlayer(msg.Player)
	if r == nil {
		r = w.voidRoom
	}
	resp := message.ExitsResponse{
		Response: message.NewSuccessfulResponse("exits"),
		ExitInfo: r.GetExitInfo(),
	}
	msg.Player.Send(resp)
}
