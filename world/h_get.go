package world

import (
	log "github.com/rs/zerolog/log"
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/gameserver"
)

func (w *World) handleGet(msg *gameserver.HandlerParameter, cmd command.Get) {
	// for now, just 'get' the first target if it is given
	// (multitarget stuff we'll deal with another time)
	if cmd.Target == "" {
		msg.Fail(event.NoTarget)
		return
	}

	room := w.getRoomContainingPlayer(msg.Player)
	items := room.Inventory.NameOrAlias(cmd.Target)
	if len(items) == 0 {
		msg.Fail(event.TargetNotFound)
		return
	}
	// for now just grab the 0th one
	item := items[0]

	if !item.Definition.Gettable() {
		msg.Fail(event.TargetNotGettable)
		return
	}

	// remove from room
	if err := room.Inventory.Remove(item); err != nil {
		// uh oh failed to remove from room
		log.Error().Err(err).Str("zone", room.Zone.Name).Str("playerName", msg.Player.Name()).Str("room", room.Name).Msgf("handle_get: error removing item %s (%s) from room", item.Id, item.Definition.ObjectId.String())
		msg.Fail(event.RemoveFromRoomError)
		return
	}
	// add to player
	msg.Player.Inventory().Add(item)

	room.Send(event.Got{
		Actor: msg.Player.Name(),
		Item:  item.Definition.ShortDescription,
	})
}
