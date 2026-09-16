package spaces

import (
	"fmt"
	"iter"
	"slices"
	"uuid"

	"github.com/trasa/watchmud/mobile"
)

type RoomMobs struct {
	byId  map[uuid.UUID]*mobile.Instance
	order []*mobile.Instance
}

func NewRoomMobs() *RoomMobs {
	return &RoomMobs{
		byId: make(map[uuid.UUID]*mobile.Instance),
	}
}

func (rm *RoomMobs) All() iter.Seq[*mobile.Instance] {
	return slices.Values(rm.order)
}

func (rm *RoomMobs) Find(target string) (*mobile.Instance, bool) {
	for _, inst := range rm.order {
		if inst.Matches(target) {
			return inst, true
		}
	}
	return nil, false
}

func (rm *RoomMobs) Remove(instance *mobile.Instance) error {
	if _, exists := rm.byId[instance.Id()]; !exists {
		return fmt.Errorf("remove: mob instance %s named %s does not exist in room", instance.Id(), instance.Name())
	}
	delete(rm.byId, instance.Id())
	if i := slices.IndexFunc(rm.order, func(m *mobile.Instance) bool {
		return m.Id() == instance.Id()
	}); i >= 0 {
		rm.order = slices.Delete(rm.order, i, i+1)
	}
	return nil
}

func (rm *RoomMobs) Add(inst *mobile.Instance) error {
	if _, exists := rm.byId[inst.Id()]; exists {
		return fmt.Errorf("add: mob instance %s named %s is already in the room", inst.Id(), inst.Name())
	}
	rm.byId[inst.Id()] = inst
	rm.order = append(rm.order, inst)
	return nil
}
