package world

import (
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/message"
	"log"
)

func (w *World) handleGet(msg *gameserver.HandlerParameter) {
	// for now, just 'get' the first target if it is given
	// (multitarget stuff we'll deal with another time)
	getreq := msg.Message.GetGetRequest()

	if len(getreq.Targets) == 0 {
		msg.Player.Send(message.GetResponse{Success: false, ResultCode: "NO_TARGET"})
		return
	}

	room := w.getRoomContainingPlayer(msg.Player)
	if instPtr, ok := room.GetInventoryByName(getreq.Targets[0]); ok {
		// target is in room

		if !instPtr.IsGettable() {
			msg.Player.Send(message.GetResponse{
				Success: false,
				ResultCode: "NO_TAKE",
			})
			return
		}

		// add to player
		if err := msg.Player.AddInventory(instPtr); err != nil {
			// uh oh failed to add
			log.Printf("Get: Error while getting, Player %s adding Inventory %v: %s",
				msg.Player.GetName(), instPtr, err)
			msg.Player.Send(message.GetResponse{Success: false, ResultCode: "ADD_INVENTORY_ERROR"})
			return
		}

		// remove from room
		if err := room.RemoveInventory(instPtr); err != nil {
			// uh oh failed to remove from room
			log.Printf("Get: Error while removing from room: Player %s Inventory %s: %s", msg.Player.GetName(), instPtr.Id(), err)
			msg.Player.RemoveInventory(instPtr)
			msg.Player.Send(message.GetResponse{Success: false, ResultCode: "REMOVE_FROM_ROOM_ERROR"})
			return
		}
		// success!
		msg.Player.Send(message.GetResponse{Success: true, ResultCode: "OK"})

		// tell everyone else in room too
		room.SendExcept(msg.Player,
			message.GetNotification{
				Success:    true,
				ResultCode: "OK",
				Target:     instPtr.Definition.Name, // what should be sent?! needs to handle various articles, plural...
				PlayerName: msg.Player.GetName(),
			})

		return
	} else {
		// nothing here with that name
		msg.Player.Send(message.GetResponse{Success: false, ResultCode: "TARGET_NOT_FOUND"})
	}
}
