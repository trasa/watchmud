package world

import (
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/object"
)

func (w *World) handleInventory(msg *message.IncomingMessage) {
	items := []message.InventoryDescription{}
	for _, instPtr := range msg.Player.GetInventoryMap() {
		items = append(items, ObjectInstanceToDescription(*instPtr.(*object.Instance)))
	}
	resp := message.InventoryResponse{
		Response:       message.NewSuccessfulResponse("inv"),
		InventoryItems: items,
	}
	msg.Player.Send(resp)
}

func ObjectInstanceToDescription(item object.Instance) message.InventoryDescription {
	return message.InventoryDescription{
		Id:               item.InstanceId,
		ShortDescription: item.Definition.ShortDescription,
		ObjectCategories: item.Definition.Categories.ToList(),
	}
}
