package player

import (
	"github.com/trasa/watchmud/object"
)

type InventoryItem struct {
	Id               string
	ShortDescription string
	ObjectCategory   object.Category
	Quantity         int
}
