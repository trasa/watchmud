package spaces

import (
	"fmt"
	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/object"
)

type Zone struct {
	Id                string
	Rooms             map[string]*Room
	ObjectDefinitions map[string]*object.Definition // id -> object.Definition
	MobileDefinitions map[string]*mobile.Definition // id -> mobile.Definition
	Name              string
}

func NewZone(id string, name string) *Zone {
	return &Zone{
		Id:                id,
		Name:              name,
		Rooms:             make(map[string]*Room),
		ObjectDefinitions: make(map[string]*object.Definition),
		MobileDefinitions: make(map[string]*mobile.Definition),
	}
}

func (z *Zone) AddRoom(r *Room) {
	r.Zone = z
	z.Rooms[r.Id] = r
}

func (z *Zone) AddObjectDefinition(obj *object.Definition) {
	obj.ZoneId = z.Id
	z.ObjectDefinitions[obj.Id] = obj
}

func (z *Zone) AddMobileDefinition(mob *mobile.Definition) {
	mob.ZoneId = z.Id
	z.MobileDefinitions[mob.Id] = mob
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
