package world

import (
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/message"
)

func (w *World) handleShowEquipment(msg *gameserver.HandlerParameter) {
	var items []*message.ShowEquipmentResponse_EquipmentInfo
	for loc, inst := range msg.Player.Slots().GetAll() {
		items = append(items, &message.ShowEquipmentResponse_EquipmentInfo{
			Id:               inst.Id(),
			ShortDescription: inst.Definition.ShortDescription,
			SlotLocation:     int32(loc),
		})
	}
	msg.Player.Send(message.ShowEquipmentResponse{
		Success:       true,
		ResultCode:    "OK",
		EquipmentInfo: items,
	})
}
