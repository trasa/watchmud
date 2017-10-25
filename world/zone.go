package world

import "fmt"

type Zone struct {
	Id    string
	Rooms map[string]*Room
	Name  string
}

func NewZone(id string, name string) *Zone {
	return &Zone{
		Id:    id,
		Name:  name,
		Rooms: make(map[string]*Room),
	}
}

func (z *Zone) String() string {
	return fmt.Sprintf("(Zone %s: '%s')", z.Id, z.Name)
}

func (z *Zone) DoZoneActivity() error {

	// Time to do a zone reset?
	// Reset Modes:
	// 0 - Never Reset
	// 1 - Reset only when no players are present in the zone
	// 2 - Reset even if players are present
	// TODO http://www.circlemud.org/cdp/building/building-6.html
	return nil
}

func (z *Zone) AddRoom(r *Room) {
	z.Rooms[r.Id] = r
}
