package world

import (
	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/slot"
)

func (w *World) handleEquip(msg *gameserver.HandlerParameter, cmd command.Equip) {
	if cmd.Slot <= slot.None {
		msg.Fail(event.NoSlotGiven)
		return
	}
	if cmd.Target == "" {
		msg.Fail(event.NoTarget)
		return
	}

	target, err := parseTarget(cmd.Target)
	if err != nil {
		log.Debug().Err(err).Str("player", msg.Player.Name()).Str("target", cmd.Target).Msg("equip: can't parse target")
		msg.Fail(event.ParseError)
		return
	}

	// TODO need to sort this by priority of how we count "2.x"
	// and make sense of other info in the target data structure
	// TODO other soring things
	objectsToEquip := msg.Player.Inventory().GetByNameOrAlias(target.Name)
	if len(objectsToEquip) == 0 {
		// you don't have one
		msg.Fail(event.TargetNotFound)
		return
	}

	// TODO for now, use the first one returned
	objectToEquip := objectsToEquip[0]

	// do you already have something equipped in that location?
	if msg.Player.Slots().Get(cmd.Slot) != nil {
		msg.Fail(event.LocationInUse)
		return
	}

	// can this object be equiped there?
	if cmd.Slot != objectToEquip.Definition.WearLocation {
		msg.Fail(event.CantWearThere)
		return
	}
	// success
	msg.Player.Slots().Set(cmd.Slot, objectToEquip)
	msg.Player.Send(event.Equipped{})
}
