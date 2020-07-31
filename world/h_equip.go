package world

import (
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud-message/slot"
	"github.com/trasa/watchmud/gameserver"
)

func (w *World) handleEquip(msg *gameserver.HandlerParameter) {
	equipReq := msg.Message.GetEquipRequest()
	requestedLocation := slot.Location(equipReq.SlotLocation)

	if equipReq.SlotLocation <= 0 {
		_ = msg.Player.Send(message.EquipResponse{
			Success:    false,
			ResultCode: "NO_SLOT_GIVEN",
		})
		return
	}
	if equipReq.Target == "" {
		_ = msg.Player.Send(message.EquipResponse{
			Success:    false,
			ResultCode: "NO_TARGET",
		})
		return
	}

	target, err := parseTarget(equipReq.Target)
	if err != nil {
		_ = msg.Player.Send(message.EquipResponse{Success: false, ResultCode: "PARSE_ERROR_" + err.Error()})
		return
	}

	// TODO need to sort this by priority of how we count "2.x"
	// and make sense of other info in the target data structure
	// TODO other soring things
	objectsToEquip := msg.Player.Inventory().GetByNameOrAlias(target.Name)
	if len(objectsToEquip) == 0 {
		// you don't have one
		_ = msg.Player.Send(message.EquipResponse{
			Success:    false,
			ResultCode: "TARGET_NOT_FOUND",
		})
		return
	}

	// TODO for now, use the first one returned
	objectToEquip := objectsToEquip[0]

	// do you already have something equipped in that location?
	if msg.Player.Slots().Get(requestedLocation) != nil {
		_ = msg.Player.Send(message.EquipResponse{
			Success:    false,
			ResultCode: "LOCATION_IN_USE",
		})
		return
	}

	// can this object be equiped there?
	if requestedLocation != objectToEquip.Definition.WearLocation {
		_ = msg.Player.Send(message.EquipResponse{
			Success:    false,
			ResultCode: "CANT_WEAR_THERE",
		})
		return
	}
	// success
	msg.Player.Slots().Set(requestedLocation, objectToEquip)
	_ = msg.Player.Send(message.EquipResponse{
		Success:    true,
		ResultCode: "OK",
	})
}
