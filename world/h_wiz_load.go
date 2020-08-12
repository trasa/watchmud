package world

import (
	message "github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/spaces"
)

func (w *World) handleLoad(msg *gameserver.HandlerParameter) {
	// TODO figure out the level of the user and if they are allowed to run this wizcommand!
	loadRequest := msg.Message.GetLoadRequest()

	targetRoom := w.getRoomContainingPlayer(msg.Player)
	if targetRoom == nil {
		msg.Player.Send(message.LoadResponse{Success: false, ResultCode: "YOU_ARE_NOT_IN_A_ROOM"})
		return
	}

	if loadRequest.Zone == "" {
		loadRequest.Zone = targetRoom.Zone.Id
	}
	logWizCommand(msg.Player, "load",
		"Player %s is creating %s of %s.%s", msg.Player.GetName(), loadRequest.Type, loadRequest.Zone, loadRequest.Id)

	if loadRequest.Type == "mob" {
		w.handleLoadCreateMob(msg, loadRequest, targetRoom)
	} else if loadRequest.Type == "obj" {
		w.handleLoadCreateObject(msg, loadRequest, targetRoom)
	} else {
		_ = msg.Player.Send(message.LoadResponse{Success: false, ResultCode: "UNKNOWN_TYPE"})
	}
}

func (w *World) handleLoadCreateMob(msg *gameserver.HandlerParameter, request *message.LoadRequest, targetRoom *spaces.Room) {
	// get the zone we're looking for a mob in
	z := w.GetZone(request.Zone)
	if z == nil {
		_ = msg.Player.Send(message.LoadResponse{Success: false, ResultCode: "UNKNOWN_ZONE"})
		return
	}

	// get the definition of this mob from that zone
	mobDefn := z.MobileDefinitions[request.Id]
	if mobDefn == nil {
		_ = msg.Player.Send(message.LoadResponse{Success: false, ResultCode: "UNKNOWN_ID"})
		return
	}

	// create instance of the mob
	inst := mobile.NewInstance(mobDefn)

	// add instance to room via the world
	// have to do it this way so that the World has appropriate bookkeeping,
	// if you add directly to the target room then you'll cause problems.
	w.AddMobile(inst, targetRoom)

	// success
	_ = msg.Player.Send(message.LoadResponse{Success: true, ResultCode: "OK"})
}

func (w *World) handleLoadCreateObject(msg *gameserver.HandlerParameter, request *message.LoadRequest, targetRoom *spaces.Room) {
	// get the zone we're looking for an instance in
	z := w.GetZone(request.Zone)
	if z == nil {
		_ = msg.Player.Send(message.LoadResponse{Success: false, ResultCode: "UNKNOWN_ZONE"})
		return
	}

	// get the definition of this object from that zone
	objDefn := z.ObjectDefinitions[request.Id]
	if objDefn == nil {
		_ = msg.Player.Send(message.LoadResponse{Success: false, ResultCode: "UNKNOWN_ID"})
		return
	}

	// create instance of the item
	inst, err := object.NewInstance(objDefn)
	if err != nil {
		_ = msg.Player.Send(message.LoadResponse{Success: false, ResultCode: "CREATE_INSTANCE_FAILED"})
		return
	}

	// add instance to room
	if err := targetRoom.AddInventory(inst); err != nil {
		_ = msg.Player.Send(message.LoadResponse{Success: false, ResultCode: "ADD_ROOM_INVENTORY_FAILED"})
		return
	}
	// success
	_ = msg.Player.Send(message.LoadResponse{Success: true, ResultCode: "OK"})
}
