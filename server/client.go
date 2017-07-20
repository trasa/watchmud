package server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

type Client interface {
	Send(message interface{}) // todo return err
	SetPlayer(player *Player)
	GetPlayer() *Player
	Close()
}

type WebClient struct {
	conn   *websocket.Conn  // websocket connection
	source chan interface{} // sends down to client
	quit   chan interface{} // used to terminate clients
	Player *Player
}

func newWebClient(c *websocket.Conn) *WebClient {
	return &WebClient{
		conn:   c,
		source: make(chan interface{}, 256),
		quit:   make(chan interface{}),
	}
}

func (c *WebClient) Send(message interface{}) {
	c.source <- message
}

func (c *WebClient) GetPlayer() *Player {
	return c.Player
}

func (c *WebClient) SetPlayer(player *Player) {
	log.Print("setting player!")
	c.Player = player
}

func (c *WebClient) Close() {
	close(c.quit)
}

func (c *WebClient) readPump() {
	defer c.conn.Close()

	for {
		body := make(map[string]string)
		err := c.conn.ReadJSON(&body)
		if err != nil {
			log.Printf("read error: %s", err)
			// TODO terminate /disconnect player
			return
		}
		log.Printf("message body: %s", body)
		GameServerInstance.incomingMessageBuffer <- newIncomingMessage(c, body)
	}
}

func (c *WebClient) writePump() {
	defer c.conn.Close()

	c.source = make(chan interface{}, 10)
	for {
		select {
		case message := <-c.source:
			log.Printf("writing %s", message)
			err := c.conn.WriteJSON(message)
			if err != nil {
				log.Printf("Write Error: %v", err)
				// TODO terminate/disconnect player
				return
			}
		case <-c.quit:
			return // terminate the client
		}
	}
}

func (c *WebClient) String() string {
	return fmt.Sprintf("(WebClient conn: %v, Player %s", c.conn != nil, c.Player)
}
