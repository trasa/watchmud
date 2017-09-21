package object

import (
	"errors"
	"fmt"
)

// The Instances of the Definitions in the world around you.
// That ShinySword in your hand has certain properties, some of
// which were inherited by what it means to be a ShinySword (ObjectType)
// and others which have happened to that particular instance
// (soul-bound to you, made invisible, with some damage to the hilt).
type Instance struct {
	// TODO what properties does an Instance have?
	InstanceId string
	Definition *Definition
}

func NewInstance(instanceId string, defn *Definition) *Instance {
	return &Instance{
		InstanceId: instanceId,
		Definition: defn,
	}
}

// map InstanceIDs to Instance objects
type InstanceMap map[string]*Instance

func (m InstanceMap) Add(inst *Instance) error {
	if _, exists := m[inst.InstanceId]; exists {
		return errors.New(fmt.Sprintf("instance %s already exists in map", inst.InstanceId))
	} else {
		m[inst.InstanceId] = inst
	}
	return nil
}
