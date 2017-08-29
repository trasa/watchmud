package message

import (
	"github.com/trasa/watchmud/direction"
)

type MoveRequest struct {
	Request
	Direction direction.Direction `json:"direction"`
}

type MoveResponse struct {
	Response
	RoomDescription
}

type EnterRoomNotification struct {
	Response
	PlayerName string `json:"player"`
}

type LeaveRoomNotification struct {
	Response
	PlayerName string              `json:"player"`
	Direction  direction.Direction `json:"direction"`
}
