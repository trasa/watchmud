package world

import (
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud/gameserver"
	"log"
)

func (w *World) initializeHandlerMap() {
	w.handlerMap = map[string]func(parameter *gameserver.HandlerParameter){
		"GameMessage_DropRequest":          w.handleDrop,
		"GameMessage_EquipRequest":         w.handleEquip,
		"GameMessage_ExitsRequest":         w.handleExits,
		"GameMessage_GetRequest":           w.handleGet,
		"GameMessage_InventoryRequest":     w.handleInventory,
		"GameMessage_KillRequest":          w.handleKill,
		"GameMessage_LogoutRequest":        w.handleLogout,
		"GameMessage_LookRequest":          w.handleLook,
		"GameMessage_MoveRequest":          w.handleMove,
		"GameMessage_PingRequest":          w.handlePing,
		"GameMessage_RoomStatusRequest":    w.handleRoomStatus,
		"GameMessage_SayRequest":           w.handleSay,
		"GameMessage_ShowEquipmentRequest": w.handleShowEquipment,
		"GameMessage_TellRequest":          w.handleTell,
		"GameMessage_TellAllRequest":       w.handleTellAll,
		"GameMessage_WearRequest":          w.handleWear,
		"GameMessage_WhoRequest":           w.handleWho,
	}
	return
}

func (w *World) HandleIncomingMessage(msg *gameserver.HandlerParameter) {
	handler := w.handlerMap[message.DecodeTypeName(msg.Message.Inner)]
	if handler == nil {
		log.Printf("world.HandleIncomingMessage: UNHANDLED messageType: %v, body %s", msg.Message.Inner, msg.Message)
		msg.Client.Send(message.ErrorResponse{
			Success:    false,
			ResultCode: "UNKNOWN_MESSAGE_TYPE",
		})
	} else {
		msg.Player.ResetDirtyFlag()
		handler(msg)
		// if the player object has changed, persist the changes to the database
		if msg.Player.IsDirty() {
			// todo save
		}
	}
}
