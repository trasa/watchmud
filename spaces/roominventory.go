package spaces

import (
	"fmt"
	"uuid"

	"github.com/trasa/watchmud/object"
)

type RoomInventory struct {
	byInstanceId   map[uuid.UUID]*object.Instance
	insertionOrder []*object.Instance
}

func NewRoomInventory() *RoomInventory {
	return &RoomInventory{
		byInstanceId: make(map[uuid.UUID]*object.Instance),
	}
}

// GetAll instances in the room sorted by their insertion order
func (ri *RoomInventory) GetAll() []*object.Instance {
	return ri.insertionOrder
}

// InstanceId finds the instance with this id in the room
func (ri *RoomInventory) InstanceId(id uuid.UUID) (*object.Instance, bool) {
	if o, ok := ri.byInstanceId[id]; ok {
		return o, true
	}
	return nil, false
}

// Name finds the instances with this name in the room.
// TODO this needs to become much more sophisticated...
// Note that there is much left undone by this implementation
// (stacks of items, aliases...)
func (ri *RoomInventory) Name(name string) []*object.Instance {
	var result []*object.Instance
	for _, inst := range ri.GetAll() {
		if inst.Definition.Name == name {
			result = append(result, inst)
		}
	}
	return result
}

func (ri *RoomInventory) NameOrAlias(target string) []*object.Instance {
	var result []*object.Instance
	for _, inst := range ri.GetAll() {
		if inst.Definition.Name == target || inst.Definition.HasAlias(target) {
			result = append(result, inst)
		}
	}
	return result
}

func (ri *RoomInventory) Add(inst *object.Instance) error {
	if _, exists := ri.InstanceId(inst.Id); exists {
		return fmt.Errorf("instance id %s already exists in room inventory", inst.Id)
	}
	ri.byInstanceId[inst.Id] = inst
	ri.insertionOrder = append(ri.insertionOrder, inst)
	return nil
}

func (ri *RoomInventory) Remove(inst *object.Instance) error {
	if _, exists := ri.InstanceId(inst.Id); !exists {
		return fmt.Errorf("instance id %s does not exist in room inventory", inst.Id)
	}
	delete(ri.byInstanceId, inst.Id)

	for i, o := range ri.insertionOrder {
		if o.Id == inst.Id {
			ri.insertionOrder = append(ri.insertionOrder[:i], ri.insertionOrder[i+1:]...)
			break
		}
	}
	return nil
}
