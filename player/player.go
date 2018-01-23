package player

import (
	"github.com/satori/go.uuid"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/thing"
)

// see https://play.golang.org/p/zPLyr3ZOM0 (first attempt)
// then see https://play.golang.org/p/z5athD5fV3 (client is an interface, but now pointer woes)
//noinspection GoNameStartsWithPackageName
type Player interface {
	Send(innerMessage interface{}) error
	GetName() string
	AddInventory(instance *object.Instance) error
	RemoveInventory(instance *object.Instance) error
	GetInventoryByName(definitionId string) (*object.Instance, bool)
	GetInventoryById(id uuid.UUID) (*object.Instance, bool)
	GetAllInventory() []*object.Instance
	Inventory() thing.Map
	Slots() *object.Slots
}
