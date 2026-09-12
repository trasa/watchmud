package world

import (
	"uuid"

	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/spaces"
)

func (w *World) handleLoad(msg *gameserver.HandlerParameter, cmd command.Load) {
	// TODO figure out the level of the user and if they are allowed to run this wizcommand!
	targetRoom := w.getRoomContainingPlayer(msg.Player)
	if targetRoom == nil {
		msg.Fail(event.YouAreNotInARoom)
		return
	}

	if cmd.Zone == "" {
		cmd.Zone = targetRoom.Zone.Id
	}
	logWizCommand(msg.Player, "load",
		"Player %s is creating %s of %s.%s", msg.Player.Name(), cmd.Type, cmd.Zone, cmd.Id)

	switch cmd.Type {
	case "mob":
		w.handleLoadCreateMob(msg, cmd, targetRoom)
	case "obj":
		w.handleLoadCreateObject(msg, cmd, targetRoom)
	default:
		msg.Fail(event.UnknownType)
	}
}

func (w *World) handleLoadCreateMob(msg *gameserver.HandlerParameter, cmd command.Load, targetRoom *spaces.Room) {
	// get the zone we're looking for a mob in
	z := w.Zone(cmd.Zone)
	if z == nil {
		msg.Fail(event.UnknownZone)
		return
	}

	// get the definition of this mob from that zone
	mobDefn := z.MobileDefinitions[cmd.Id]
	if mobDefn == nil {
		msg.Fail(event.UnknownId)
		return
	}

	// create instance of the mob
	inst := mobile.NewInstance(mobDefn)

	// add instance to room via the world
	// have to do it this way so that the World has appropriate bookkeeping,
	// if you add directly to the target room then you'll cause problems.
	w.AddMobile(inst, targetRoom)

	msg.Player.Send(event.Loaded{})
}

func (w *World) handleLoadCreateObject(msg *gameserver.HandlerParameter, cmd command.Load, targetRoom *spaces.Room) {
	// get the zone we're looking for an instance in
	z := w.Zone(cmd.Zone)
	if z == nil {
		msg.Fail(event.UnknownZone)
		return
	}

	// get the definition of this object from that zone
	definition := z.ObjectDefinitions[cmd.Id]
	if definition == nil {
		msg.Fail(event.UnknownDefinitionId)
		return
	}

	// create instance of the item
	inst := object.NewInstance(uuid.New(), definition)

	// add instance to room
	if err := targetRoom.Inventory.Add(inst); err != nil {
		msg.Fail(event.AddRoomInventoryFailed)
		return
	}
	msg.Player.Send(event.Loaded{})
}
