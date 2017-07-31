package web

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/trasa/watchmud/gameserver"
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/player"
	"log"
)

var gameServerInstance gameserver.Instance

func Init(gs gameserver.Instance) {
	gameServerInstance = gs
}

type Client struct {
	conn   *websocket.Conn  // websocket connection
	source chan interface{} // sends down to client
	quit   chan interface{} // used to terminate clients
	Player player.Player
}

func newClient(c *websocket.Conn) *Client {
	return &Client{
		conn:   c,
		source: make(chan interface{}, 256),
		quit:   make(chan interface{}),
	}
}

func (c *Client) Send(message interface{}) {
	c.source <- message
}

func (c *Client) GetPlayer() player.Player {
	return c.Player
}

func (c *Client) SetPlayer(player player.Player) {
	c.Player = player
}

func (c *Client) Close() {
	close(c.quit)
}

func (c *Client) readPump() {
	defer c.conn.Close()

	for {
		body := make(map[string]string)
		if err := c.conn.ReadJSON(&body); err != nil {
			log.Printf("read error: %s", err)
			// TODO terminate / disconnect player
			return
		}

		var request message.Request
		var err error
		if body["format"] == "line" {
			request, err = translateLineToRequest(body["value"])
		} else {
			request, err = translateToRequest(body)
		}
		if err != nil {
			log.Printf("translation error: %s", err)
			return // TODO terminate / disconnect player
		}
		gameServerInstance.Receive(message.New(c, request))
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	c.source = make(chan interface{}, 10)
	for {
		select {
		case msg := <-c.source:
			err := c.conn.WriteJSON(msg)
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

func (c *Client) String() string {
	return fmt.Sprintf("(WebClient conn: %v, Player %s", c.conn != nil, c.Player)
}
