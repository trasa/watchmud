package message

import (
	"github.com/trasa/watchmud/direction"
)

type MoveRequest struct {
	Request
	Direction direction.Direction
}

type MoveResponse struct {
	Response
	RoomDescription RoomDescription
}

type EnterRoomNotification struct {
	Response
	Name string
}

type LeaveRoomNotification struct {
	Response
	Name      string
	Direction direction.Direction
}
