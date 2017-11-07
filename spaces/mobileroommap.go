package spaces

import (
	"github.com/trasa/syncmap"
	"github.com/trasa/watchmud/mobile"
	"sync"
)

type mobileToRoom map[*mobile.Instance]*Room

type MobileRoomMap struct {
	sync.RWMutex
	mobileToRoom  mobileToRoom
	roomToMobiles syncmap.MapList
}

func NewMobileRoomMap() *MobileRoomMap {
	return &MobileRoomMap{
		mobileToRoom:  make(mobileToRoom),
		roomToMobiles: syncmap.NewMapList(),
	}
}

func (m *MobileRoomMap) GetRoomForMobile(mob *mobile.Instance) *Room {
	m.RLock()
	defer m.RUnlock()
	return m.mobileToRoom[mob]
}

func (m *MobileRoomMap) GetAllMobiles() (mobs []*mobile.Instance) {
	m.RLock()
	defer m.RUnlock()
	for m := range m.mobileToRoom {
		mobs = append(mobs, m)
	}
	return
}

func (m *MobileRoomMap) Add(mob *mobile.Instance, r *Room) {
	m.Lock()
	defer m.Unlock()
	m.mobileToRoom[mob] = r
	m.roomToMobiles.Add(r, mob)
	r.AddMobile(mob)
}

func (m *MobileRoomMap) Remove(mob *mobile.Instance) {
	m.Lock()
	defer m.Unlock()
	r := m.mobileToRoom[mob]
	delete(m.mobileToRoom, mob)
	if r != nil {
		m.roomToMobiles.RemoveItem(r, mob)
		r.RemoveMobile(mob)
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
