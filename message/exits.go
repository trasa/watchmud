package message

import "github.com/trasa/watchmud/direction"

type ExitsRequest struct {
	Request
}

type ExitsResponse struct {
	Response
	// maps directionAbbreviation to name of room
	ExitInfo []DirectionToRoomName
}

type DirectionToRoomName struct {
	Direction direction.Direction
	RoomName  string
}
