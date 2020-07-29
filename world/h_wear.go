package world

import (
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/gameserver"
)

func (w *World) handleWear(msg *gameserver.HandlerParameter) {
	objectsToWear := msg.Player.Inventory().GetByNameOrAlias(msg.Message.GetWearRequest().Target)

	if len(objectsToWear) == 0 {
		// nothing in inventory with that name
		_ = msg.Player.Send(message.WearResponse{Success: false, ResultCode: "TARGET_NOT_FOUND"})
		return
	}

	// TODO for now only using first object returned
	objectToWear := objectsToWear[0]

	if !objectToWear.Definition.CanWear() {
		_ = msg.Player.Send(message.WearResponse{Success: false, ResultCode: "CANT_WEAR_THAT"})
		return
	}

	// figure out what the wear location is:
	//		was one provided? (for now, instances can only be worn in one place)
	// 		so we ignore the given location
	loc := objectToWear.Definition.WearLocation

	// is something else already in the location?
	if msg.Player.Slots().IsSlotInUse(loc) {
		_ = msg.Player.Send(message.WearResponse{Success: false, ResultCode: "IN_USE"})
		return
	}

	// otherwise add the item to the location
	// TODO fix this so that Set() only takes one thing?
	msg.Player.Slots().Set(objectToWear.Definition.WearLocation, objectToWear)
	_ = msg.Player.Send(message.WearResponse{Success: true, ResultCode: "OK"})
}
