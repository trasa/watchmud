package world

import (
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/gameserver"
)

func (w *World) handleRestore(msg *gameserver.HandlerParameter, cmd command.Restore) {
	// TODO figure out the level of the user and if they are allowed to run this wizcommand!
	targetRoom := w.getRoomContainingPlayer(msg.Player)
	if targetRoom == nil {
		msg.Fail(event.YouAreNotInARoom)
		return
	}

	logWizCommand(msg.Player, "restore", "Player %s is attempting to restore %s",
		msg.Player.Name(), cmd.Target)

	// find a matching player
	if targetPlayer, found := targetRoom.FindPlayer(cmd.Target); found {
		// TODO implement restore
		//targetPlayer.Restore()
		targetRoom.Notify(event.Restored{
			IsPlayer: true,
			Target:   targetPlayer.Name(),
		})
		return
	}

	// find a matching mob
	if targetMob, found := targetRoom.FindMobile(cmd.Target); found {
		targetMob.Restore()
		targetRoom.Notify(event.Restored{
			IsPlayer: false,
			Target:   targetMob.Name(),
		})
		return
	}
	msg.Fail(event.TargetNotFound)
}
