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
	Notification
	PlayerName string `json:"player"`
}

type LeaveRoomNotification struct {
	Notification
	PlayerName string              `json:"player"`
	Direction  direction.Direction `json:"direction"`
}
