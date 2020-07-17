package object

import (
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/trasa/watchmud/behavior"
	"github.com/rs/zerolog/log"
)

// The Instances of the Definitions in the world around you.
// That ShinySword in your hand has certain properties, some of
// which were inherited by what it means to be a ShinySword (Definition)
// and others which have happened to that particular instance
// (soul-bound to you, made invisible, with some damage to the hilt).
type Instance struct {
	InstanceId uuid.UUID
	Definition *Definition
}

func (i *Instance) Id() string {
	return i.InstanceId.String()
}
func (i *Instance) CanEquipWeapon() bool {
	return i.Definition.CanEquipWeapon()
}

func (i *Instance) IsGettable() bool {
	return !i.Definition.Behaviors.Contains(behavior.NoTake)
}

func NewInstance(defn *Definition) (*Instance, error) {
	id := uuid.NewV4()
	return NewExistingInstance(id, defn)
}

func NewExistingInstance(id uuid.UUID, defn *Definition) (inst *Instance, err error) {
	if defn == nil {
		log.Error().Msgf("Error: asked to create instance for id %s with null definition!", id)
		return nil, errors.New("Tried to create instance with null definition")
	}
	return &Instance{
		InstanceId: id,
		Definition: defn,
	}, nil
}
