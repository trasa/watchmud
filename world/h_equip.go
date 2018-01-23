package world

import (
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/slot"
	"reflect"
)

func (w *World) handleEquip(msg *gameserver.HandlerParameter) {
	equipReq := msg.Message.GetEquipRequest()
	if equipReq.SlotLocation <= 0 {
		msg.Player.Send(message.EquipResponse{
			Success:    false,
			ResultCode: "NO_SLOT_GIVEN",
		})
		return
	}
	if equipReq.Target == "" {
		msg.Player.Send(message.EquipResponse{
			Success:    false,
			ResultCode: "NO_TARGET",
		})
		return
	}

	instPtr, ok := msg.Player.GetInventoryByName(equipReq.Target)
	if !ok {
		// you don't have one
		msg.Player.Send(message.EquipResponse{
			Success:    false,
			ResultCode: "TARGET_NOT_FOUND",
		})
		return
	}
	// player has target
	err := msg.Player.Slots().Set(slot.Location(equipReq.SlotLocation), instPtr)
	if err != nil {
		msg.Player.Send(message.EquipResponse{
			Success:    false,
			ResultCode: reflect.TypeOf(err).Name(),
		})
		return
	}

	msg.Player.Send(message.EquipResponse{
		Success:    true,
		ResultCode: "OK",
	})
}
