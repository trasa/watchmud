package world

import (
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/message"
	"log"
)

func (w *World) handleKill(msg *gameserver.HandlerParameter) {
	// TODO make different kill command for killing a player vs
	// killing a mob. for now this will just be killing mobs.

	killRequest := msg.Message.GetKillRequest()

	// if you're already in a fight, you can't start a new fight
	if msg.Player.IsFighting() {
		msg.Player.Send(message.KillResponse{Success: false, ResultCode: "ALREADY_FIGHTING"})
		return
	}
	// figure out if the target of your fight is valid
	//	are they in the room (still)
	//	are they something you are allowed to fight (no_fight, other flags... objects...)
	//  does this room allow fighting..
	//  and so on

	room := w.getRoomContainingPlayer(msg.Player)
	mobileInstance, exists := room.FindMobile(killRequest.Target)
	log.Printf("found mobile %v, exists %v", mobileInstance, exists)

	// begin a fight with that target (or join an existing fight if there's
	// already one going on with that target)

	// TODO
	msg.Player.Send(message.KillResponse{Success: false, ResultCode: "TODO"})

}
