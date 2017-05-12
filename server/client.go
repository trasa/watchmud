package server

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
		GameServerInstance.incomingMessageBuffer <- newMessage(c, body)
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
