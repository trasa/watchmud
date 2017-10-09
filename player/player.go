package player

import (
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/thing"
)

// see https://play.golang.org/p/zPLyr3ZOM0 (first attempt)
// then see https://play.golang.org/p/z5athD5fV3 (client is an interface, but now pointer woes)
//noinspection GoNameStartsWithPackageName
type Player interface {
	Send(message interface{}) // todo return err
	GetName() string
	GetInventoryMap() thing.Map
	AddInventory(instance *object.Instance) error
	RemoveInventory(instance *object.Instance) error
}
