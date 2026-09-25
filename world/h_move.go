package world

import (
	"github.com/watchmud/watchmud/command"
	"github.com/watchmud/watchmud/event"
	"github.com/watchmud/watchmud/gameserver"
)

func (w *World) handleMove(msg *gameserver.HandlerParameter, cmd command.Move) {
	// is player in a fight?
	if w.fightLedger.InFight(msg.Player) {
		msg.Fail(event.InAFight)
		return
	}

	playerRoom := w.getPlayerRoom(msg.Player)
	dir := cmd.Direction

	msg.Player.Log().Trace().Msgf("player wants to move %s", dir.String())

	// can player go in that direction?
	targetRoom := playerRoom.DestinationRoom(dir)
	if targetRoom == nil {
		msg.Fail(event.CantGoThatWay)
		return
	}
	w.movePlayer(msg.Player, dir, playerRoom, targetRoom)
	msg.Player.Send(targetRoom.DescriptionExcept(msg.Player))
}
