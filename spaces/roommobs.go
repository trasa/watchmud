package spaces

import (
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
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
	if _, exists := rm.byId[inst.InstanceId]; !exists {
		return errors.New(fmt.Sprintf("instance id %s does not exist in room", inst.InstanceId))
	}
	delete(rm.byId, inst.InstanceId)
	// TODO other indexes
	return nil
}

func (rm *RoomMobs) Add(inst *mobile.Instance) error {
	if _, exists := rm.byId[inst.InstanceId]; exists {
		return errors.New(fmt.Sprintf("instance id %s is already in the room", inst.InstanceId))
	}
	rm.byId[inst.InstanceId] = inst
	// TODO other indexes
	return nil
}
