package world

import "github.com/trasa/watchmud/message"

func (w *World) handleGet(msg *message.IncomingMessage) {
	// for now, just 'get' the first target if it is given
	// (multitarget stuff we'll deal with another time)
	getreq := msg.Request.(message.GetRequest)
	if len(getreq.Targets) == 0 {
		msg.Player.Send(message.GetResponse{
			Response: message.NewUnsuccessfulResponse("get", "NO_TARGET"),
		})
	} else {
		// TODO
		msg.Player.Send(message.GetResponse{
			Response: message.NewUnsuccessfulResponse("get", "NOT_IMPLEMENTED_YET"),
		})
		//room := w.getRoomContainingPlayer(msg.Player)
		//room.Objects
	}
}
