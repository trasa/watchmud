package world

import (
	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/direction"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/gameserver"
)

const fleeAttempts = 6

func (w *World) handleFlee(msg *gameserver.HandlerParameter, cmd command.Flee) {

	if !w.fightLedger.IsFighting(msg.Player) {
		// you're not in a fight?!
		msg.Fail(event.NoFight)
		return
	}

	// TODO stunned, disabled, other status effects
	room := w.getRoomContainingPlayer(msg.Player)

	for range fleeAttempts {
		room.Send(event.Fleeing{Who: msg.Player.Name()})

		// pick a direction at random out of all possible
		i, err := w.roller.IntN(len(direction.All))
		if err != nil {
			log.Error().Err(err).Msg("flee: failed to generate random direction")
			msg.Fail(event.CantFlee)
			return
		}
		// is there an exit?
		// TODO need to check other things like if the door is locked, etc.
		dir := direction.All[i]
		if room.HasExit(dir) {
			// success!
			room.Send(event.Fled{Who: msg.Player.Name()})
			w.movePlayer(msg.Player, dir, room, room.DestinationRoom(dir))
			w.fightLedger.EndAllFightsWith(msg.Player.Id())
			return
		}
		// you can't escape
		room.Send(event.FleeAttemptFailed{Who: msg.Player.Name()})
	}
	msg.Fail(event.CantFlee)
}
