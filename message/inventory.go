package message

import (
	"github.com/trasa/watchmud/object"
)

type InventoryRequest struct {
	Request
}

type InventoryResponse struct {
	Response
	InventoryItems []InventoryDescription
}

type InventoryDescription struct {
	Id               string
	ShortDescription string
	ObjectCategory   object.Category
	Quantity         int
}
