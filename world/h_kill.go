package world

import (
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/spaces"
)

func (w *World) handleKill(msg *gameserver.HandlerParameter, cmd command.Kill) {
	// TODO make different kill command for killing a player vs
	// killing a mob. for now this will just be killing mobs.

	// if you're already in a fight, you can't start a new fight
	if w.fightLedger.IsFighting(msg.Player) {
		msg.Fail(event.AlreadyFighting)
		return
	}

	// figure out if the target of your fight is valid
	//	are they in the room (still)
	room := w.getRoomContainingPlayer(msg.Player)
	mobileInstance, exists := room.FindMobile(cmd.Target)
	if !exists {
		msg.Fail(event.TargetNotFound)
		return
	}

	//  does this room allow fighting..
	if room.Flag(spaces.RoomFlagNoFight) {
		msg.Fail(event.NoFightRoom)
		return
	}

	//	are they something you are allowed to fight (no_fight, other flags... objects...)
	if mobileInstance.Definition.HasFlag(mobile.FlagNoFight) {
		msg.Fail(event.NoFight)
		return
	}

	// begin a fight with that target (or join an existing fight if there's
	// already one going on with that target)
	// TODO reimplement
	/*
		if err := w.fightLedger.Fight(msg.Player, mobileInstance, room.Zone.Id, room.Id); err != nil {
			log.Error().Err(err).Msg("kill: couldn't start the fight")
			msg.Fail(event.InternalError)
			return
		}
	*/
	msg.Player.Send(event.Attacking{Target: mobileInstance.Name()})
}
