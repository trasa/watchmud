package server

import (
	"testing"
)

func TestGameServer_handleMessage_unknownMessageType(t *testing.T) {
	server := newGameServer()
	body := map[string]string{
		"msg_type": "asdfasdf",
	}
	server.handleMessage(newMessage(&Client{}, body))
}

func TestGameServer_handleLogin(t *testing.T) {
	// what to do here...
	server := newGameServer()
	body := map[string]string{
		"msg_type":    "login",
		"player_name": "testdood",
		"password":    "password",
	}
	server.handleLogin(newMessage(&Client{}, body))
	// TODO verify state of world after
}
