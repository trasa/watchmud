package world

import (
	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/gameserver"
)

func (w *World) handleMove(msg *gameserver.HandlerParameter, cmd command.Move) {

	// is player in a fight?
	// TODO reimplement
	/*
		if w.fightLedger.IsFighting(msg.Player) || w.fightLedger.IsBeingFought(msg.Player) {
			msg.Fail(event.InAFight)
			return
		}
	*/
	playerRoom := w.getRoomContainingPlayer(msg.Player)
	dir := cmd.Direction

	log.Trace().Str("player", msg.Player.Name()).Str("room", playerRoom.Name).Msgf("player wants to move %s", dir.String())

	// can player go in that direction?
	targetRoom := playerRoom.DestinationRoom(dir)
	if targetRoom == nil {
		msg.Fail(event.CantGoThatWay)
		return
	}
	w.movePlayer(msg.Player, dir, playerRoom, targetRoom)
	msg.Player.Send(targetRoom.DescriptionExcept(msg.Player))
}
