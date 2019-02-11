package player

import "fmt"

type Location struct {
	ZoneId string
	RoomId string
}

func NewLocation(zoneId *string, roomId *string) *Location {
	if zoneId == nil || roomId == nil {
		return &Location{}
	}
	return &Location{
		ZoneId: *zoneId,
		RoomId: *roomId,
	}
}

func (l Location) String() string {
	return fmt.Sprintf("(%s - %s)", l.ZoneId, l.RoomId)
}
