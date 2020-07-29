package player

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/satori/go.uuid"
	"github.com/trasa/watchmud/object"
)

//noinspection GoNameStartsWithPackageName
type Inventory struct {
	byId         map[uuid.UUID]*object.Instance // instance_id -> instance obj
	byDefinition map[string][]*object.Instance  // zone.definition_id -> list of instances

	added   map[uuid.UUID]*object.Instance
	removed map[uuid.UUID]*object.Instance
}

func NewInventory() *Inventory {
	return &Inventory{
		byId:         make(map[uuid.UUID]*object.Instance),
		byDefinition: make(map[string][]*object.Instance),
		added:        make(map[uuid.UUID]*object.Instance),
		removed:      make(map[uuid.UUID]*object.Instance),
	}
}

func (pi *Inventory) GetAll() (result []*object.Instance) {
	for _, inst := range pi.byId {
		result = append(result, inst)
	}
	return
}

func (pi *Inventory) GetByInstanceId(id uuid.UUID) (inst *object.Instance, exists bool) {
	inst, exists = pi.byId[id]
	return
}

func (pi *Inventory) GetByName(name string) (inst *object.Instance, exists bool) {
	for _, inst := range pi.GetAll() {
		if inst.Definition.Name == name {
			return inst, true
		}
	}
	return nil, false
}

// Find the objects in the inventory that match this name or alias
func (pi *Inventory) GetByNameOrAlias(target string) (objects []*object.Instance) {
	// TODO handle case where target is "2.knife" (return the 2nd knife)
	// TODO handle all the other target cases
	objects = []*object.Instance{}
	for _, obj := range pi.GetAll() {
		log.Debug().Msgf("consider %+v", obj)
		if obj.Definition.Name == target || obj.Definition.HasAlias(target) {
			objects = append(objects, obj)
		}
	}
	return objects
}

func (pi *Inventory) findPosition(inst *object.Instance) int {
	for pos, i := range pi.byDefinition[inst.Definition.Id()] {
		if i.InstanceId == inst.InstanceId {
			return pos
		}
	}
	return -1
}

// Add an object into the inventory. This marks the inventory as dirty and records the change.
func (pi *Inventory) Add(inst *object.Instance) error {
	err := pi.Load(inst)
	if err != nil {
		return err
	}
	pi.added[inst.InstanceId] = inst
	return nil
}

// Load an object into the inventory without marking the inventory as dirty or changed.
func (pi *Inventory) Load(inst *object.Instance) error {
	if _, exists := pi.byId[inst.InstanceId]; exists {
		return errors.New(fmt.Sprintf("instance id %s already exists in player inventory", inst.InstanceId))
	}
	pi.byId[inst.InstanceId] = inst
	pi.byDefinition[inst.Definition.Id()] = append(pi.byDefinition[inst.Definition.Id()], inst)
	return nil
}

// Remove an object from the inventory. This marks the inventory as dirty and records the change.
func (pi *Inventory) Remove(inst *object.Instance) error {
	if _, exists := pi.byId[inst.InstanceId]; !exists {
		return errors.New(fmt.Sprintf("instance id %s does not exist in player inventory", inst.InstanceId))
	}
	delete(pi.byId, inst.InstanceId)
	pos := pi.findPosition(inst)
	pi.byDefinition[inst.Definition.Id()] = append(pi.byDefinition[inst.Definition.Id()][:pos], pi.byDefinition[inst.Definition.Id()][pos+1:]...)
	pi.removed[inst.InstanceId] = inst
	return nil
}

// Get the objects that were added since the last time
func (pi *Inventory) GetAdded() (result []*object.Instance) {
	for _, inst := range pi.added {
		result = append(result, inst)
	}
	return
}

// get the objects that were removed since hte last time
func (pi *Inventory) GetRemoved() (result []*object.Instance) {
	for _, inst := range pi.removed {
		result = append(result, inst)
	}
	return
}

func (pi *Inventory) IsDirty() bool {
	return len(pi.added) > 0 || len(pi.removed) > 0
}

func (pi *Inventory) ResetDirtyFlag() {
	pi.added = make(map[uuid.UUID]*object.Instance)
	pi.removed = make(map[uuid.UUID]*object.Instance)
}
