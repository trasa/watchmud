package world

import (
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/spaces"
)

func (w *World) handleKill(msg *gameserver.HandlerParameter) {
	// TODO make different kill command for killing a player vs
	// killing a mob. for now this will just be killing mobs.
	killRequest := msg.Message.GetKillRequest()

	// if you're already in a fight, you can't start a new fight
	if w.fightLedger.IsFighting(msg.Player) {
		msg.Player.Send(message.KillResponse{Success: false, ResultCode: "ALREADY_FIGHTING"})
		return
	}

	// figure out if the target of your fight is valid
	//	are they in the room (still)
	room := w.getRoomContainingPlayer(msg.Player)
	mobileInstance, exists := room.FindMobile(killRequest.Target)
	if !exists {
		msg.Player.Send(message.KillResponse{Success: false, ResultCode: "TARGET_NOT_FOUND"})
		return
	}

	//  does this room allow fighting..
	if room.HasFlag(spaces.RoomFlagNoFight) {
		msg.Player.Send(message.KillResponse{Success: false, ResultCode: "NO_FIGHT_ROOM"})
		return
	}

	//	are they something you are allowed to fight (no_fight, other flags... objects...)
	if mobileInstance.Definition.HasFlag(mobile.FlagNoFight) {
		msg.Player.Send(message.KillResponse{Success: false, ResultCode: "NO_FIGHT"})
		return
	}

	// begin a fight with that target (or join an existing fight if there's
	// already one going on with that target)

	// TODO
	msg.Player.Send(message.KillResponse{Success: false, ResultCode: "TODO"})

}
