package world

import (
	log "github.com/rs/zerolog/log"
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/gameserver"
)

func (w *World) handleGet(msg *gameserver.HandlerParameter) {
	// for now, just 'get' the first target if it is given
	// (multitarget stuff we'll deal with another time)
	getreq := msg.Message.GetGetRequest()
	// TODO findMode not supported yet
	findMode := message.FindMode(getreq.FindMode)
	if findMode == message.FindIndividual && len(getreq.Target) == 0 {
		msg.Player.Send(message.GetResponse{Success: false, ResultCode: "NO_TARGET"})
		return
	}

	room := w.getRoomContainingPlayer(msg.Player)
	items := room.Inventory.NameOrAlias(getreq.Target)
	if len(items) == 0 {
		msg.Player.Send(message.GetResponse{Success: false, ResultCode: "TARGET_NOT_FOUND"})
		return
	}
	// for now just grab the 0th one
	item := items[0]

	if !item.Definition.Gettable() {
		msg.Player.Send(message.GetResponse{
			Success:    false,
			ResultCode: "TARGET_NOT_GETTABLE",
		})
		return
	}

	// remove from room
	if err := room.Inventory.Remove(item); err != nil {
		// uh oh failed to remove from room
		log.Error().Err(err).Str("zone", room.Zone.Name).Str("playerName", msg.Player.Name).Str("room", room.Name).Msgf("handle_get: error removing item %s (%s) from room", item.Id, item.Definition.ObjectId.String())
		msg.Player.Send(message.GetResponse{Success: false, ResultCode: "REMOVE_FROM_ROOM_ERROR"})
		return
	}
	// add to player
	msg.Player.Inventory().Add(item)
	msg.Player.Send(message.GetResponse{Success: true, ResultCode: "OK"})

	// tell everyone else in room too
	room.SendExcept(msg.Player,
		message.GetNotification{
			Success:    true,
			ResultCode: "OK",
			Target:     item.Definition.Name,
			PlayerName: msg.Player.Name,
		})
	return
}
