package world

import (
	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/gameserver"
)

func (w *World) handleDrop(msg *gameserver.HandlerParameter, cmd command.Drop) {
	if cmd.Target == "" {
		msg.Fail(event.NoTarget)
		return
	}

	target, err := parseTarget(cmd.Target)
	if err != nil {
		// the parse error is for us, not for the player: it's things like
		// TOO_MANY_DOTS and strconv's complaint about "x.knife".
		log.Debug().Err(err).Str("player", msg.Player.Name()).Str("target", cmd.Target).Msg("drop: can't parse target")
		msg.Fail(event.ParseError)
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
		msg.Fail(event.TargetNotFound)
		return
	}

	// TODO for now, using the first item returned
	// player has target
	objectToDrop := objectsToDrop[0]

	// is the object cursed?
	// TODO cursed

	// is the object being held or otherwise in use?
	if msg.Player.Slots().IsItemInUse(objectToDrop) {
		msg.Fail(event.TargetInUse)
		return
	}

	// add to room
	if err := room.Inventory.Add(objectToDrop); err != nil {
		// failed to add to room..
		log.Error().Err(err).Str("player", msg.Player.Name()).Stringer("id", objectToDrop.Id).Msg("drop: error while adding to room")
		msg.Fail(event.AddToRoomError)
		return
	}

	// remove from player
	msg.Player.Inventory().Remove(objectToDrop)

	// one event, both audiences: the renderer says "Dropped." to the actor
	// and "bob drops a knife." to everyone else.
	room.Send(event.Dropped{
		Actor: msg.Player.Name(),
		Item:  objectToDrop.Definition.ShortDescription, // rendered to clients, so use "a knife"
	})
}
