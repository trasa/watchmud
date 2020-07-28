package world

import (
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/gameserver"
)

func (w *World) handleWear(msg *gameserver.HandlerParameter) {
	wearreq := msg.Message.GetWearRequest()

	if instPtr, ok := msg.Player.Inventory().GetByNameOrAlias(wearreq.Target); ok {
		if !instPtr.Definition.CanWear() {
			msg.Player.Send(message.WearResponse{Success: false, ResultCode: "CANT_WEAR_THAT"})
			return
		}

		// figure out what the wear location is:
		//		was one provided? (for now, instances can only be worn in one place)
		// 		so we ignore the given location
		loc := instPtr.Definition.WearLocation

		// is something else already in the location?
		if msg.Player.Slots().IsSlotInUse(loc) {
			msg.Player.Send(message.WearResponse{Success: false, ResultCode: "IN_USE"})
			return
		}

		// otherwise add the item to the location
		msg.Player.Slots().Set(instPtr.Definition.WearLocation, instPtr)
		msg.Player.Send(message.WearResponse{Success: true, ResultCode: "OK"})

	} else {
		// nothing in inventory with that name
		msg.Player.Send(message.WearResponse{Success: false, ResultCode: "TARGET_NOT_FOUND"})
	}
}
