package spaces

import (
	"fmt"
	"uuid"

	"github.com/trasa/watchmud/mobile"
)

type RoomMobs struct {
	byId map[uuid.UUID]*mobile.Instance
}

func NewRoomMobs() *RoomMobs {
	return &RoomMobs{
		byId: make(map[uuid.UUID]*mobile.Instance),
	}
}

func (rm *RoomMobs) GetAll() (result []*mobile.Instance) {
	for _, inst := range rm.byId {
		result = append(result, inst)
	}
	return result
}

func (rm *RoomMobs) Find(target string) (inst *mobile.Instance, exists bool) {
	for _, inst := range rm.GetAll() {
		if inst.Definition.Name == target {
			return inst, true
		}
		if inst.Definition.HasAlias(target) {
			return inst, true
		}
	}
	return nil, false
}

func (rm *RoomMobs) Remove(inst *mobile.Instance) error {
	if _, exists := rm.byId[inst.Id()]; !exists {
		return fmt.Errorf("remove: mob instance %s named %s does not exist in room", inst.Id(), inst.Name())
	}
	delete(rm.byId, inst.Id())
	return nil
}

func (rm *RoomMobs) Add(inst *mobile.Instance) error {
	if _, exists := rm.byId[inst.Id()]; exists {
		return fmt.Errorf("add: mob instance %s named %s is already in the room", inst.Id(), inst.Name())
	}
	rm.byId[inst.Id()] = inst
	return nil
}
