package world

import (
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/message"
)

func (w *World) handleKill(msg *gameserver.HandlerParameter) {
	// TODO
	msg.Player.Send(message.KillResponse{Success: false, ResultCode: "TODO"})
}
