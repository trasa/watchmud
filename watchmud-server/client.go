package main

import (
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	// websocket connection
	conn *websocket.Conn
	// buffered channel of outbound messages
	send          chan []byte // what to send?
	authenticated bool
}

func newClient(c *websocket.Conn) *Client {
	return &Client{
		conn: c,
		send: make(chan []byte, 256),
	}
}

func (c *Client) readPump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		body := make(map[string]string)
		err := c.conn.ReadJSON(&body)
		if err != nil {
			log.Printf("read error: %s", err)
			break
		}

		log.Printf("message body: %s", body)
		GameServerInstance.incomingMessageBuffer <- newMessage(c, &body)
		// need to put this message into a queue of messages to be handled
		/*
			if !c.authenticated {
				// must authenticate first, only allowable message
				if !c.authenticate(message) {
					// TODO some sort of useful response here
					log.Printf("Failed to authenticate %s", message)
				} else {
					// TODO do login (create player, land into a room, and so on)
					p := world.NewPlayer(message["player_name"], message["player_name"])
					//p.Room = world.World.StartRoom
					c.Player = p
				}
			} else {
				// TODO do someting useful with message
			}
		*/
		//c.send <- []byte("abcdefg")
	}
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case message := <-c.send:
			log.Printf("writing %s", message)
			c.conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func (c *Client) authenticate(message map[string]string) bool {
	// TODO authenticate stuff..
	c.authenticated = true
	return true
}
