package spaces

import (
	"github.com/watchmud/watchmud/mobile"
)

type mobileToRoom map[*mobile.Instance]*Room

// MobileRoomMap is which room each mob is in. The other direction, which mobs
// a room holds, is the room's own list; see ROADMAP "Dual location
// bookkeeping".
type MobileRoomMap struct {
	mobileToRoom mobileToRoom
}

func NewMobileRoomMap() *MobileRoomMap {
	return &MobileRoomMap{
		mobileToRoom: make(mobileToRoom),
	}
}

func (m *MobileRoomMap) GetRoomForMobile(mob *mobile.Instance) *Room {
	return m.mobileToRoom[mob]
}

func (m *MobileRoomMap) GetAllMobiles() (mobs []*mobile.Instance) {
	for m := range m.mobileToRoom {
		mobs = append(mobs, m)
	}
	return
}

func (m *MobileRoomMap) Add(mob *mobile.Instance, r *Room) {
	m.mobileToRoom[mob] = r
	// TODO error handling
	r.AddMobile(mob)
}

func (m *MobileRoomMap) Remove(mob *mobile.Instance) {
	r := m.mobileToRoom[mob]
	delete(m.mobileToRoom, mob)
	if r != nil {
		r.RemoveMobile(mob) // TODO error handling
	}
}

func (m *MobileRoomMap) GetMobileDefinitionCount(mobileDefinitionId string) (count int) {
	for mob := range m.mobileToRoom {
		if mob.Definition.Id == mobileDefinitionId {
			count++
		}
	}
	return
}
