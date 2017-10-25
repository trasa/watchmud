package world

import (
	"github.com/trasa/syncmap"
	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/spaces"
	"sync"
)

type mobileToRoom map[*mobile.Instance]*spaces.Room

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

func (m *MobileRoomMap) Get(mob *mobile.Instance) *spaces.Room {
	m.RLock()
	defer m.RUnlock()
	return m.mobileToRoom[mob]
}

func (m *MobileRoomMap) Add(mob *mobile.Instance, r *spaces.Room) {
	m.Lock()
	defer m.Unlock()
	m.mobileToRoom[mob] = r
	m.roomToMobiles.Add(r, mob)
}

func (m *MobileRoomMap) Remove(mob *mobile.Instance) {
	m.Lock()
	defer m.Unlock()
	r := m.mobileToRoom[mob]
	delete(m.mobileToRoom, mob)
	if r != nil {
		m.roomToMobiles.RemoveItem(r, mob)
		r.Mobs.Remove(mob)
	}
}
