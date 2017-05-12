package main

import "log"

type Message struct {
	Client *Client
	Body   *map[string]string
}

type GameServer struct {
	incomingMessageBuffer chan *Message
}

func newMessage(client *Client, body *map[string]string) *Message {
	return &Message{
		Client: client,
		Body:   body,
	}
}

func newGameServer() *GameServer {
	return &GameServer{
		incomingMessageBuffer: make(chan *Message),
	}
}

func (server *GameServer) run() {
	for {
		select {
		case message := <-server.incomingMessageBuffer:
			log.Printf("server processed incoming message: %s", message.Body)
		}
	}
}
