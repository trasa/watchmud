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

	// TODO handle "coins" (target.Quantity)

	room := w.getPlayerRoom(msg.Player)
	objectsToDrop := targetsIn(target, msg.Player.Inventory().All())
	if len(objectsToDrop) == 0 {
		msg.Fail(event.TargetNotFound)
		return
	}

	dropped := 0
	for _, objectToDrop := range objectsToDrop {
		// is the object cursed?
		// TODO cursed

		// is the object being held or otherwise in use? naming one thing you
		// have on is worth saying so; "drop all" while wearing armor is not.
		if msg.Player.Equipment().ItemEquipped(objectToDrop) {
			if !target.All {
				msg.Fail(event.TargetInUse)
				return
			}
			continue
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
		dropped++
	}

	// everything named was something you have on
	if dropped == 0 {
		msg.Fail(event.TargetInUse)
	}
}
