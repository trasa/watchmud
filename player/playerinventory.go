package player

import (
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/trasa/watchmud/object"
)

//noinspection GoNameStartsWithPackageName
type PlayerInventory struct {
	byId         map[uuid.UUID]*object.Instance // instance_id -> instance obj
	byDefinition map[string][]*object.Instance  // zone.definition_id -> list of instances
}

func NewPlayerInventory() *PlayerInventory {
	return &PlayerInventory{
		byId:         make(map[uuid.UUID]*object.Instance),
		byDefinition: make(map[string][]*object.Instance),
	}
}

func (pi *PlayerInventory) GetAll() (result []*object.Instance) {
	for _, inst := range pi.byId {
		result = append(result, inst)
	}
	return
}

func (pi *PlayerInventory) GetByInstanceId(id uuid.UUID) (inst *object.Instance, exists bool) {
	inst, exists = pi.byId[id]
	return
}

func (pi *PlayerInventory) GetByName(name string) (inst *object.Instance, exists bool) {
	for _, inst := range pi.GetAll() {
		if inst.Definition.Name == name {
			return inst, true
		}
	}
	return nil, false
}

func (pi *PlayerInventory) GetByNameOrAlias(target string) (inst *object.Instance, exists bool) {
	for _, inst := range pi.GetAll() {
		if inst.Definition.Name == target {
			return inst, true
		}
		if inst.Definition.HasAlias(target) {
			return inst, true
		}
	}
	return nil, false
}

func (pi *PlayerInventory) Find(target string) (*object.Instance, bool) {
	return pi.GetByNameOrAlias(target)
}

func (pi *PlayerInventory) findPosition(inst *object.Instance) int {
	for pos, i := range pi.byDefinition[inst.Definition.Id()] {
		if i.InstanceId == inst.InstanceId {
			return pos
		}
	}
	return -1
}

func (pi *PlayerInventory) Add(inst *object.Instance) error {
	if _, exists := pi.byId[inst.InstanceId]; exists {
		return errors.New(fmt.Sprintf("instance id %s already exists in player inventory", inst.InstanceId))
	}
	pi.byId[inst.InstanceId] = inst
	pi.byDefinition[inst.Definition.Id()] = append(pi.byDefinition[inst.Definition.Id()], inst)
	return nil
}

func (pi *PlayerInventory) Remove(inst *object.Instance) error {
	if _, exists := pi.byId[inst.InstanceId]; !exists {
		return errors.New(fmt.Sprintf("instance id %s does not exist in player inventory", inst.InstanceId))
	}
	delete(pi.byId, inst.InstanceId)
	pos := pi.findPosition(inst)
	pi.byDefinition[inst.Definition.Id()] = append(pi.byDefinition[inst.Definition.Id()][:pos], pi.byDefinition[inst.Definition.Id()][pos+1:]...)
	return nil
}
