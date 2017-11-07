package spaces

import (
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/trasa/watchmud/object"
)

type RoomInventory struct {
	byId         map[uuid.UUID]*object.Instance // instance_id -> instance object
	byDefinition map[string][]*object.Instance  // (zone)(definition_id) -> list of Instances
}

func NewRoomInventory() *RoomInventory {
	return &RoomInventory{
		byId:         make(map[uuid.UUID]*object.Instance),
		byDefinition: make(map[string][]*object.Instance),
	}
}

func (ri *RoomInventory) GetAll() (result []*object.Instance) {
	for _, inst := range ri.byId {
		result = append(result, inst)
	}
	return
}

func (ri *RoomInventory) GetByInstanceId(id uuid.UUID) (inst *object.Instance, exists bool) {
	inst, exists = ri.byId[id]
	return
}

// Find the instance with this name in the room
// note this needs to become much more sophisticated...
// Note that there is much left undone by this implementation
// (stacks of items, aliases...)
func (ri *RoomInventory) GetByName(name string) (inst *object.Instance, exists bool) {
	for _, inst := range ri.GetAll() {
		if inst.Definition.Name == name {
			return inst, true
		}
	}
	return nil, false
}

func (ri *RoomInventory) Add(inst *object.Instance) (err error) {
	if _, exists := ri.byId[inst.InstanceId]; exists {
		return errors.New(fmt.Sprintf("instance id %s already exists in room_inventory", inst.InstanceId))
	}

	ri.byId[inst.InstanceId] = inst
	ri.byDefinition[inst.Definition.Id()] = append(ri.byDefinition[inst.Definition.Id()], inst)
	return nil
}

func (ri *RoomInventory) Remove(inst *object.Instance) (err error) {
	if _, exists := ri.byId[inst.InstanceId]; !exists {
		return errors.New(fmt.Sprintf("instance id %s does not exist in room_inventory", inst.InstanceId))
	}
	delete(ri.byId, inst.InstanceId)
	pos := ri.findPosition(inst)
	ri.byDefinition[inst.Definition.Id()] = append(ri.byDefinition[inst.Definition.Id()][:pos], ri.byDefinition[inst.Definition.Id()][pos+1:]...)
	return nil
}

func (ri *RoomInventory) findPosition(inst *object.Instance) int {
	for pos, i := range ri.byDefinition[inst.Definition.Id()] {
		if i.InstanceId == inst.InstanceId {
			return pos
		}
	}
	return -1
}
