package world

import (
	"github.com/trasa/watchmud/message"
	"log"
)

func (w *World) handleDrop(msg *message.IncomingMessage) {
	dropReq := msg.Request.(message.DropRequest)
	if dropReq.Target == "" {
		msg.Player.Send(message.DropResponse{
			Response: message.NewUnsuccessfulResponse("drop", "NO_TARGET"),
		})
		return
	}
	room := w.getRoomContainingPlayer(msg.Player)
	if instPtr, ok := msg.Player.GetInventoryMap()[dropReq.Target]; ok {
		// player has target

		// add to room
		if err := room.Inventory.Add(instPtr); err != nil {
			// failed to add to room..
			log.Printf("Drop: Error while adding to room, player %s id %s; %s",
				msg.Player.GetName(),
				instPtr.InstanceId,
				err)
			msg.Player.Send(message.DropResponse{
				Response: message.NewUnsuccessfulResponse("drop", "ADD_TO_ROOM_ERROR"),
			})
			return
		}

		// remove from player
		if err := msg.Player.RemoveInventory(instPtr); err != nil {
			// failed to remove from player
			log.Printf("Drop: error while removing from player: %s id %s; %s",
				msg.Player.GetName(),
				instPtr.InstanceId,
				err)

			room.Inventory.Remove(instPtr)
			msg.Player.Send(message.DropResponse{
				Response: message.NewUnsuccessfulResponse("drop", "REMOVE_FROM_PLAYER_ERROR"),
			})
			return
		}

		// success
		msg.Player.Send(message.DropResponse{
			Response: message.NewSuccessfulResponse("drop"),
		})
	} else {
		// not found
		msg.Player.Send(message.DropResponse{
			Response: message.NewUnsuccessfulResponse("drop", "TARGET_NOT_FOUND"),
		})
	}
}
