package world

import (
	log "github.com/rs/zerolog/log"
	"github.com/watchmud/watchmud/command"
	"github.com/watchmud/watchmud/event"
	"github.com/watchmud/watchmud/gameserver"
)

func (w *World) handleGet(msg *gameserver.HandlerParameter, cmd command.Get) {
	if cmd.Target == "" {
		msg.Fail(event.NoTarget)
		return
	}
	if cmd.From != "" {
		w.getFrom(msg, cmd)
		return
	}

	target, err := parseTarget(cmd.Target)
	if err != nil {
		// the parse error is for us, not for the player: it's things like
		// TOO_MANY_DOTS and strconv's complaint about "x.knife".
		log.Debug().Err(err).Str("player", msg.Player.Name()).Str("target", cmd.Target).Msg("get: can't parse target")
		msg.Fail(event.ParseError)
		return
	}

	// TODO handle "coins" (target.Quantity)

	room := w.getPlayerRoom(msg.Player)
	items := targetsIn(target, room.Inventory.All())
	if len(items) == 0 {
		msg.Fail(event.TargetNotFound)
		return
	}

	got := 0
	for _, item := range items {
		if !item.Definition.Gettable() {
			// naming one thing that can't be picked up is worth saying so;
			// "get all" in a room with a fountain in it is not.
			if !target.All {
				msg.Fail(event.TargetNotGettable)
				return
			}
			continue
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
		got++
	}

	// everything named was something you can't pick up
	if got == 0 {
		msg.Fail(event.TargetNotGettable)
	}
}
