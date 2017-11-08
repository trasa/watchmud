package world

import (
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/message"
	"log"
)

func (w *World) handleDrop(msg *gameserver.HandlerParameter) {
	dropReq := msg.Message.GetDropRequest()
	if dropReq.Target == "" {
		msg.Player.Send(message.DropResponse{
			Success: false, ResultCode: "NO_TARGET",
		})
		return
	}
	room := w.getRoomContainingPlayer(msg.Player)
	if instPtr, ok := msg.Player.GetInventoryByName(dropReq.Target); ok {
		// player has target

		// add to room
		if err := room.AddInventory(instPtr); err != nil {
			// failed to add to room..
			log.Printf("Drop: Error while adding to room, player %s id %s; %s",
				msg.Player.GetName(),
				instPtr.Id(),
				err)
			msg.Player.Send(message.DropResponse{
				Success: false, ResultCode: "ADD_TO_ROOM_ERROR",
			})
			return
		}

		// remove from player
		if err := msg.Player.RemoveInventory(instPtr); err != nil {
			// failed to remove from player
			log.Printf("Drop: error while removing from player: %s id %s; %s",
				msg.Player.GetName(),
				instPtr.Id(),
				err)

			room.RemoveInventory(instPtr)
			msg.Player.Send(message.DropResponse{
				Success: false, ResultCode: "REMOVE_FROM_PLAYER_ERROR",
			})
			return
		}

		// success
		msg.Player.Send(message.DropResponse{
			Success: true, ResultCode: "OK",
		})
		// tell everybody about it
		room.SendExcept(msg.Player,
			message.DropNotification{
				Success:    true,
				ResultCode: "OK",
				PlayerName: msg.Player.GetName(),
				Target:     instPtr.Definition.Name, // what should this be?! "knife", "a knife", "those knives" ...
			})
	} else {
		// not found
		msg.Player.Send(message.DropResponse{
			Success: false, ResultCode: "TARGET_NOT_FOUND",
		})
	}
}
