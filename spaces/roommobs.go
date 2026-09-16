package spaces

import (
	"fmt"
	"iter"
	"uuid"

	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/ordered"
)

// RoomMobs is the mobs standing in a room, in the order they arrived.
type RoomMobs struct {
	mobs *ordered.List[uuid.UUID, *mobile.Instance]
}

func NewRoomMobs() *RoomMobs {
	return &RoomMobs{
		mobs: ordered.NewList[uuid.UUID]((*mobile.Instance).Id),
	}
}

// All the mobs in the room, in the order they arrived.
func (rm *RoomMobs) All() iter.Seq[*mobile.Instance] {
	return rm.mobs.All()
}

// Find the mob this target names -- the longest-standing one, if the room
// holds more than one of them.
func (rm *RoomMobs) Find(target string) (*mobile.Instance, bool) {
	return ordered.Find(rm.mobs, target)
}

func (rm *RoomMobs) Remove(inst *mobile.Instance) error {
	if err := rm.mobs.Remove(inst); err != nil {
		return fmt.Errorf("remove mob %s named %s: %w", inst.Id(), inst.Name(), err)
	}
	return nil
}

func (rm *RoomMobs) Add(inst *mobile.Instance) error {
	if err := rm.mobs.Add(inst); err != nil {
		return fmt.Errorf("add mob %s named %s: %w", inst.Id(), inst.Name(), err)
	}
	return nil
}
