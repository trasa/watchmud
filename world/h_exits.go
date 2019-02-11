package world

import (
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/gameserver"
)

func (w *World) handleExits(msg *gameserver.HandlerParameter) {
	r := w.getRoomContainingPlayer(msg.Player)
	if r == nil {
		r = w.VoidRoom
	}
	// convert directions to strings because json
	messageExitInfo := []*message.ExitInfo{}
	for _, rexit := range r.GetRoomExits(false) {
		messageExitInfo = append(messageExitInfo,
			&message.ExitInfo{
				Direction: int32(rexit.Direction),
				RoomName:  rexit.Room.Name,
			})
	}
	resp := message.ExitsResponse{
		Success:    true,
		ResultCode: "OK",
		ExitInfo:   messageExitInfo,
	}
	msg.Player.Send(resp)
}
