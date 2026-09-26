package world

import (
	"github.com/watchmud/watchmud/command"
	"github.com/watchmud/watchmud/event"
	"github.com/watchmud/watchmud/gameserver"
	"github.com/watchmud/watchmud/rules"
)

// handleConsider compares the player's power to a mob's, so they can tell
// before a fight starts whether they'll finish it. Mobs only: considering a
// player is tabled (LEVELS.md).
func (w *World) handleConsider(msg *gameserver.HandlerParameter, cmd command.Consider) {
	if cmd.Target == "" {
		msg.Fail(event.NoTarget)
		return
	}
	mob, exists := w.playerRoom(msg.Player).FindMobile(cmd.Target)
	if !exists {
		msg.Fail(event.TargetNotFound)
		return
	}
	yours, theirs := msg.Player.Power(), mob.Power()
	msg.Player.Send(event.Considered{
		Target:      mob.Definition.Name,
		TargetPower: theirs,
		YourPower:   yours,
		Delta:       rules.PowerDelta(yours, theirs),
	})
}
