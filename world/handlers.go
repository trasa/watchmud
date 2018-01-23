package world

import (
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/message"
	"log"
)

func (w *World) initializeHandlerMap() {
	w.handlerMap = map[string]func(parameter *gameserver.HandlerParameter){
		"GameMessage_DropRequest":      w.handleDrop,
		"GameMessage_EquipRequest": 	w.handleEquip,
		"GameMessage_ExitsRequest":     w.handleExits,
		"GameMessage_GetRequest":       w.handleGet,
		"GameMessage_InventoryRequest": w.handleInventory,
		"GameMessage_LogoutRequest":    w.handleLogout,
		"GameMessage_LookRequest":      w.handleLook,
		"GameMessage_MoveRequest":      w.handleMove,
		"GameMessage_PingRequest":      w.handlePing,
		"GameMessage_SayRequest":       w.handleSay,
		"GameMessage_TellRequest":      w.handleTell,
		"GameMessage_TellAllRequest":   w.handleTellAll,
		"GameMessage_WhoRequest":       w.handleWho,
	}
	return
}

func (w *World) HandleIncomingMessage(msg *gameserver.HandlerParameter) {
	handler := w.handlerMap[message.DecodeTypeName(msg.Message.Inner)]
	if handler == nil {
		log.Printf("world.HandleIncomingMessage: UNHANDLED messageType: %v, body %s", msg.Message.Inner, msg.Message)
		msg.Player.Send(message.ErrorResponse{
			Success:    false,
			ResultCode: "UNKNOWN_MESSAGE_TYPE",
		})
	} else {
		handler(msg)
	}
}
