package message

import "github.com/trasa/watchmud/object"

type InventoryRequest struct {
	Request
}

type InventoryResponse struct {
	Response
	InventoryItems []InventoryItem
}

// TODO: other stuff that makes up your inventory display ...
type InventoryItem struct {
	Id               string
	ShortDescription string
	ObjectCategory   object.Category
}
