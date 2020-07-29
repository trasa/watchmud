package world

import (
	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/gameserver"
)

func (w *World) handleDrop(msg *gameserver.HandlerParameter) {
	dropReq := msg.Message.GetDropRequest()
	if dropReq.Target == "" {
		_ = msg.Player.Send(message.DropResponse{
			Success: false, ResultCode: "NO_TARGET",
		})
		return
	}

	target, err := parseTarget(dropReq.Target)
	if err != nil {
		_ = msg.Player.Send(message.DropResponse{Success: false, ResultCode: "PARSE_ERROR_" + err.Error()})
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
		_ = msg.Player.Send(message.DropResponse{
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
		_ = msg.Player.Send(message.DropResponse{
			Success: false, ResultCode: "TARGET_IN_USE",
		})
		return
	}

	// add to room
	if err := room.AddInventory(objectToDrop); err != nil {
		// failed to add to room..
		log.Error().Msgf("Drop: Error while adding to room, player %s id %s; %s",
			msg.Player.GetName(),
			objectToDrop.Id(),
			err)
		_ = msg.Player.Send(message.DropResponse{
			Success: false, ResultCode: "ADD_TO_ROOM_ERROR",
		})
		return
	}

	// remove from player
	if err := msg.Player.Inventory().Remove(objectToDrop); err != nil {
		// failed to remove from player
		log.Error().
			Str("command", "drop").
			Str("player", msg.Player.GetName()).
			Err(err).
			Msgf("error while removing from player - object instance %s", objectToDrop.Id())

		removeFromRoomError := room.RemoveInventory(objectToDrop)
		log.Error().
			Str("command", "drop").
			Str("player", msg.Player.GetName()).
			Err(removeFromRoomError).
			Msgf("error while removing from player, removing from room (duplicate items!) object instance %s", objectToDrop.Id())

		_ = msg.Player.Send(message.DropResponse{
			Success: false, ResultCode: "REMOVE_FROM_PLAYER_ERROR",
		})
		return
	}

	// success
	_ = msg.Player.Send(message.DropResponse{
		Success: true, ResultCode: "OK",
	})
	// tell everybody about it
	room.SendExcept(msg.Player,
		message.DropNotification{
			Success:    true,
			ResultCode: "OK",
			PlayerName: msg.Player.GetName(),
			Target:     objectToDrop.Definition.Name, // what should this be?! "knife", "a knife", "those knives" ...
		})
}
