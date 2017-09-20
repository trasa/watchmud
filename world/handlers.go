package world

import (
	"github.com/trasa/watchmud/message"
	"log"
)

func (w *World) initializeHandlerMap() {
	w.handlerMap = map[string]func(*message.IncomingMessage){
		"exits":    w.handleExits,
		"inv":      w.handleInventory,
		"logout":   w.handleLogout,
		"look":     w.handleLook,
		"move":     w.handleMove,
		"say":      w.handleSay,
		"tell":     w.handleTell,
		"tell_all": w.handleTellAll,
		"who":      w.handleWho,
	}
	return
}

func (w *World) HandleIncomingMessage(msg *message.IncomingMessage) {
	messageType := msg.Request.GetMessageType()
	handler := w.handlerMap[messageType]
	if handler == nil {
		log.Printf("world.HandleIncomingMessage: UNHANDLED messageType: %s, body %s", messageType, msg.Request)
		msg.Player.Send(message.NewUnsuccessfulResponse(messageType, "UNKNOWN_MESSAGE_TYPE"))
	} else {
		handler(msg)
	}
}
