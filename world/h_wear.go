package world

import (
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/gameserver"
)

func (w *World) handleWear(msg *gameserver.HandlerParameter, cmd command.Wear) {
	objectsToWear := msg.Player.Inventory().GetByNameOrAlias(cmd.Target)

	if len(objectsToWear) == 0 {
		// nothing in inventory with that name
		msg.Fail(event.TargetNotFound)
		return
	}

	// TODO for now only using first object returned
	objectToWear := objectsToWear[0]

	if !objectToWear.Definition.Wearable() {
		msg.Fail(event.CantWearThat)
		return
	}

	// figure out what the wear location is:
	//		was one provided? (for now, instances can only be worn in one place)
	// 		so we ignore the given location
	loc := objectToWear.Definition.WearLocation

	// is something else already in the location?
	if msg.Player.Slots().IsSlotInUse(loc) {
		msg.Fail(event.InUse)
		return
	}

	// otherwise add the item to the location
	// TODO fix this so that Set() only takes one thing?
	msg.Player.Slots().Set(loc, objectToWear)
	msg.Player.Send(event.Worn{})
}
