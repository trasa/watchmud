package loader

import (
	"fmt"

	"github.com/watchmud/watchmud/rules"
)

// objectDurability is what one object definition starts out able to take.
//
// The object gets the last word, the durability table has the first, and in
// between is the whole point of the table: a builder adding a leather jerkin
// writes no number at all and gets what leather is worth. An explicit zero
// means a thing that never wears out, which is why the field is a pointer --
// otherwise every object in every existing file would be claiming it.
func objectDurability(zoneName string, obj objectEntry, table rules.DurabilityTable) (int, error) {
	if obj.Durability == nil {
		return table.MaxFor(obj.ArmorType), nil
	}
	if *obj.Durability < 0 {
		return 0, fmt.Errorf("object %s/%s: negative durability %d", zoneName, obj.Id, *obj.Durability)
	}
	return *obj.Durability, nil
}
