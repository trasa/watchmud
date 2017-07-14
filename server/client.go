package server

import (
	"github.com/gorilla/websocket"
	"log"
)

// controls terminating all clients
var GlobalQuit = make(chan interface{})

// channel for sending to all clients
var Broadcaster = make(chan interface{}, 10)

func SendToAllClients(msg interface{}) {
	Broadcaster <- msg
}

func StartAllClientDispatcher() {
	go func() {
		for {
			msg := <-Broadcaster
			clients.iter(func(c *Client) {
				c.source <- msg
			})
		}
	}()
}

type Client struct {
	conn   *websocket.Conn  // websocket connection
	source chan interface{} // sends down to client
	quit   chan interface{} // used to terminate clients
	Player *Player
}

func newClient(c *websocket.Conn) *Client {
	return &Client{
		conn:   c,
		source: make(chan interface{}, 256),
		quit:   GlobalQuit,
	}
}

func (c *Client) send(message interface{}) {
	c.source <- message
}

func (c *Client) readPump() {
	defer c.conn.Close()

	for {
		body := make(map[string]string)
		err := c.conn.ReadJSON(&body)
		if err != nil {
			log.Printf("read error: %s", err)
			clients.remove(c)
			return
		}
		log.Printf("message body: %s", body)
		GameServerInstance.incomingMessageBuffer <- newIncomingMessage(c, body)
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	c.source = make(chan interface{}, 10)
	for {
		select {
		case message := <-c.source:
			log.Printf("writing %s", message)
			err := c.conn.WriteJSON(message)
			if err != nil {
				log.Printf("Write Error: %v", err)
				clients.remove(c)
				return
			}
		case <-c.quit:
			return // terminate the client
		}
	}
}
