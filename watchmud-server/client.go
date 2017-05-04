package main

import (
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	// websocket connection
	conn *websocket.Conn
	// buffered channel of outbound messages
	send chan []byte // what to send?
}

func (c *Client) readPump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		message := make(map[string]interface{})
		err := c.conn.ReadJSON(&message)
		if err != nil {
			log.Printf("read error: %s", err)
			break
		}
		//c.hub.broadcast <- message
		// TODO do someting useful with message
		log.Printf("message: %s", message)
		c.send <- []byte("abcdefg")
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
