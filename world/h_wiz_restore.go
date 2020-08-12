package world

import (
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/gameserver"
)

func (w *World) handleRestore(msg *gameserver.HandlerParameter) {
	// TODO figure out the level of the user and if they are allowed to run this wizcommand!
	restoreRequest := msg.Message.GetRestoreRequest()

	targetRoom := w.getRoomContainingPlayer(msg.Player)
	if targetRoom == nil {
		_ = msg.Player.Send(message.RestoreResponse{Success: false, ResultCode: "YOU_ARE_NOT_IN_A_ROOM"})
		return
	}

	logWizCommand(msg.Player, "restore", "Player %s is attempting to restore %s",
		msg.Player.GetName(), restoreRequest.Target)

	// find a matching player
	if targetPlayer, found := targetRoom.FindPlayer(restoreRequest.Target); found {
		targetPlayer.Restore()
		targetRoom.Notify(message.RestoreNotification{
			IsPlayer: true,
			Target:   targetPlayer.GetName(),
		})
		_ = msg.Player.Send(message.RestoreResponse{Success: true, ResultCode: "OK"})
		return
	}

	// find a matching mob
	if targetMob, found := targetRoom.FindMobile(restoreRequest.Target); found {
		targetMob.Restore()
		targetRoom.Notify(message.RestoreNotification{
			IsPlayer: false,
			Target:   targetMob.GetName(),
		})
		_ = msg.Player.Send(message.RestoreResponse{Success: true, ResultCode: "OK"})
		return
	}
	_ = msg.Player.Send(message.RestoreResponse{Success: false, ResultCode: "TARGET_NOT_FOUND"})
}
