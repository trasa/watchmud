package world

import (
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/object"
	"log"
)

func (w *World) handleGet(msg *message.IncomingMessage) {
	// for now, just 'get' the first target if it is given
	// (multitarget stuff we'll deal with another time)
	getreq := msg.Request.(message.GetRequest)

	if len(getreq.Targets) == 0 {
		msg.Player.Send(message.GetResponse{
			Response: message.NewUnsuccessfulResponse("get", "NO_TARGET"),
		})
		return
	}

	room := w.getRoomContainingPlayer(msg.Player)
	if instPtr, ok := room.Inventory[getreq.Targets[0]]; ok {
		// target is in room

		// add to player
		if err := msg.Player.AddInventory(instPtr.(*object.Instance)); err != nil {
			// uh oh failed to add
			log.Printf("Get: Error while getting, Player %s adding Inventory %s: %s",
				msg.Player.GetName(), instPtr.Id(), err)
			msg.Player.Send(message.GetResponse{
				Response: message.NewUnsuccessfulResponse("get", "ADD_INVENTORY_ERROR"),
			})
			return
		}

		// remove from room
		if err := room.Inventory.Remove(instPtr); err != nil {
			// uh oh failed to remove from room
			log.Printf("Get: Error while removing from room: Player %s Inventory %s: %s", msg.Player.GetName(), instPtr.Id(), err)
			msg.Player.RemoveInventory(instPtr.(*object.Instance))
			msg.Player.Send(message.GetResponse{
				Response: message.NewUnsuccessfulResponse("get", "REMOVE_FROM_ROOM_ERROR"),
			})
			return
		}
		// success!
		msg.Player.Send(message.GetResponse{
			Response: message.NewSuccessfulResponse("get"),
		})

		// tell everyone else in room too
		room.SendExcept(msg.Player,
			message.GetNotification{
				Response:   message.NewSuccessfulResponse("get_notification"),
				Target:     instPtr.(*object.Instance).Definition.Name, // what should be sent?! needs to handle various articles, plural...
				PlayerName: msg.Player.GetName(),
			})

		return
	} else {
		// nothing here with that name
		msg.Player.Send(message.GetResponse{
			Response: message.NewUnsuccessfulResponse("get", "TARGET_NOT_FOUND"),
		})
	}
}
