package world

import (
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/slot"
)

func (w *World) handleEquip(msg *gameserver.HandlerParameter) {
	equipReq := msg.Message.GetEquipRequest()
	requestedLocation := slot.Location(equipReq.SlotLocation)

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

	// do you already have something equipped in that location?
	if msg.Player.Slots().Get(requestedLocation) != nil {
		msg.Player.Send(message.EquipResponse{
			Success:    false,
			ResultCode: "LOCATION_IN_USE",
		})
		return
	}

	// can this object be equiped there?
	if requestedLocation != instPtr.Definition.WearLocation {
		msg.Player.Send(message.EquipResponse{
			Success:    false,
			ResultCode: "CANT_WEAR_THERE",
		})
		return
	}
	// success
	msg.Player.Slots().Set(requestedLocation, instPtr)
	msg.Player.Send(message.EquipResponse{
		Success:    true,
		ResultCode: "OK",
	})
}
