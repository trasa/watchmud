package world

import (
	"github.com/trasa/watchmud/player"
	"testing"
)

// TODO this breaks because of channels!

// create a new test world
func newTestWorld() *World {
	return &World{
		knownPlayersByName: make(map[string]player.Player),
	}
}

func TestGameServer_handleMessage_unknownMessageType(t *testing.T) {
	t.Skip("channels are broken here")
	//world := NewWorld()
	//body := map[string]string{
	//	"msg_type": "asdfasdf",
	//}
	//world.handleIncomingMessage(newIncomingMessage(&Client{}, body))
}

func TestGameServer_handleLogin(t *testing.T) {
	t.Skip("channels are broken here")
	// what to do here...
	//server := newGameServer()
	//body := map[string]string{
	//	"msg_type":    "login",
	//	"player_name": "testdood",
	//	"password":    "password",
	//}
	// TODO fixme
	//server.handleLogin(newIncomingMessage(&Client{}, body))
	// TODO verify state of world after
}
