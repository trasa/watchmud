package world

import (
	"github.com/watchmud/watchmud/command"
	"github.com/watchmud/watchmud/event"
	"github.com/watchmud/watchmud/gameserver"
)

// handleRecall takes the player back to the start room. It began as a
// debugging command, and there'll be a better one; until then it is at least
// not a way out of a fight -- a free, certain escape would make flee, which
// can fail, pointless. Refused the way walking off is.
func (w *World) handleRecall(msg *gameserver.HandlerParameter, cmd command.Recall) {
	if w.fightLedger.InFight(msg.Player) {
		msg.Fail(event.InAFight)
		return
	}
	w.movePlayerMagically(msg.Player, w.StartRoom)
	msg.Player.Send(w.StartRoom.DescriptionExcept(msg.Player))
}
