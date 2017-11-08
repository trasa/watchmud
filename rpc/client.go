package rpc

import (
	"fmt"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/player"
	"io"
	"log"
)

// An implementation of client.Client
type client struct {
	gameServerInstance gameserver.Instance
	stream             message.MudComm_SendReceiveServer
	sendQueue          chan *message.GameMessage // sends down to client
	quit               chan interface{}          // used to terminate client
	Player             player.Player
}

func newClient(stream message.MudComm_SendReceiveServer, gs gameserver.Instance) *client {
	return &client{
		gameServerInstance: gs,
		stream:             stream,
		sendQueue:          make(chan *message.GameMessage, 256),
		quit:               make(chan interface{}),
	}
}

func (c *client) Send(inner interface{}) error {
	m, err := message.NewGameMessage(inner)
	if err != nil {
		return err
	}
	c.sendQueue <- m
	return nil
}

func (c *client) Close() {
	close(c.quit)
}

func (c *client) SetPlayer(player player.Player) {
	c.Player = player
}

func (c *client) GetPlayer() player.Player {
	return c.Player
}

func (c *client) readPump() {
	for {
		gameMessage, err := c.stream.Recv()
		if err == io.EOF {
			log.Printf("EOF received")
			c.gameServerInstance.Logout(c, "EOF received")
			return
		}
		if err != nil {
			log.Printf("RPC Read Error: %v", err)
			c.gameServerInstance.Logout(c, fmt.Sprintf("Read Error: %v", err))
			return
		}
		c.gameServerInstance.Receive(gameserver.NewHandlerParameter(c, gameMessage))
	}
}

func (c *client) writePump() {
	for {
		select {
		case msg := <-c.sendQueue:
			if err := c.stream.Send(msg); err != nil {
				log.Printf("RPC Write Error: %v", err)
				c.gameServerInstance.Logout(c, fmt.Sprintf("Write Error: %v", err))
				return
			}
		case <-c.quit:
			c.gameServerInstance.Logout(c, "QUIT channel")
			return
		}
	}
}
