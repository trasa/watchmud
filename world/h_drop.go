package world

import (
	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/gameserver"
)

func (w *World) handleDrop(msg *gameserver.HandlerParameter) {
	dropReq := msg.Message.GetDropRequest()
	if dropReq.Target == "" {
		msg.Player.Send(message.DropResponse{
			Success: false, ResultCode: "NO_TARGET",
		})
		return
	}

	target, err := parseTarget(dropReq.Target)
	if err != nil {
		msg.Player.Send(message.DropResponse{Success: false, ResultCode: "PARSE_ERROR_" + err.Error()})
		return
	}

	// TODO handle "all"
	// TODO handle "coins"

	room := w.getRoomContainingPlayer(msg.Player)

	// TODO need to sort this by priority of how we count "2.x"
	// and make sense of other info in the target data structure
	objectsToDrop := msg.Player.Inventory().GetByNameOrAlias(target.Name)
	// SORT by in use

	if len(objectsToDrop) == 0 {
		// not found
		msg.Player.Send(message.DropResponse{
			Success: false, ResultCode: "TARGET_NOT_FOUND",
		})
		return
	}

	// TODO for now, using the first item returned
	// player has target
	objectToDrop := objectsToDrop[0]

	// is the object cursed?
	// TODO cursed

	// is the object being held or otherwise in use?
	if msg.Player.Slots().IsItemInUse(objectToDrop) {
		// can't drop for 'reason'
		msg.Player.Send(message.DropResponse{
			Success: false, ResultCode: "TARGET_IN_USE",
		})
		return
	}

	// add to room
	if err := room.Inventory.Add(objectToDrop); err != nil {
		// failed to add to room..
		log.Error().Msgf("Drop: Error while adding to room, player %s id %s; %s",
			msg.Player.Name(),
			objectToDrop.Id,
			err)
		msg.Player.Send(message.DropResponse{
			Success: false, ResultCode: "ADD_TO_ROOM_ERROR",
		})
		return
	}

	// remove from player
	msg.Player.Inventory().Remove(objectToDrop)
	// success
	msg.Player.Send(message.DropResponse{
		Success: true, ResultCode: "OK",
	})
	// tell everybody about it
	room.SendExcept(msg.Player,
		message.DropNotification{
			Success:    true,
			ResultCode: "OK",
			PlayerName: msg.Player.Name(),
			Target:     objectToDrop.Definition.ShortDescription, // rendered to clients, so use "a knife"
		})
}
