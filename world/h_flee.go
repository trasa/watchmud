package world

import (
	"github.com/rs/zerolog/log"
	"github.com/watchmud/watchmud/command"
	"github.com/watchmud/watchmud/event"
	"github.com/watchmud/watchmud/gameserver"
	"github.com/watchmud/watchmud/rules"
)

const fleeAttempts = 6

func (w *World) handleFlee(msg *gameserver.HandlerParameter, cmd command.Flee) {

	if !w.fightLedger.IsFighting(msg.Player) {
		// you're not in a fight?!
		msg.Fail(event.NoFight)
		return
	}

	// TODO stunned, disabled, other status effects
	room := w.getPlayerRoom(msg.Player)

	for range fleeAttempts {
		room.Send(event.Fleeing{Who: msg.Player.Name()})

		// pick a direction at random out of all possible
		// where direction.All is ([North, East, South, West, Up, Down])
		// we're asking for [0, 6)
		i, err := w.roller.IntN(len(rules.AllUsableDirections))
		if err != nil {
			log.Error().Err(err).Msg("flee: failed to generate random direction")
			msg.Fail(event.CantFlee)
			return
		}
		dir := rules.AllUsableDirections[i]
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
