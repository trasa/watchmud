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
		requestEnvelope := message.RequestEnvelope{}
		err := c.conn.ReadJSON(&requestEnvelope)
		if err != nil {
			log.Printf("read error: %s", err)
			gameServerInstance.Logout(c, fmt.Sprintf("READ_ERROR: %s", err))
			return
		}
		if requestEnvelope.Format == "" || requestEnvelope.Value == nil {
			log.Printf("Unmarshaling error: requestEnvelope format or value is nil")
			gameServerInstance.Logout(c, "UNMARSHAL_ERROR")
			return
		}

		var request message.Request
		switch requestEnvelope.Format {
		case "line":
			// the 'line input' form of a request / command:
			// "tell bob hi there"
			request, err = translateLineToRequest(requestEnvelope.Value.(string))
		case "request":
			request, err = translateToRequest(requestEnvelope.Value.(map[string]interface{}))
		default:
			log.Println("Unhandled requestEnvelope.Format: ", requestEnvelope.Format)
			gameServerInstance.Logout(c, fmt.Sprintf("UNKNOWN_FORMAT: %s", requestEnvelope.Format))
			return
		}
		if err != nil {
			log.Printf("translation error: %s", err)
			// can't send this ourselves (not reentrant?)

			request = message.ErrorRequest{
				Request: message.RequestBase{
					MessageType: "error",
				},
				Error: err,
			}
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
				gameServerInstance.Logout(c, fmt.Sprintf("WRITE_ERROR: %s", err))
				return
			}
		case <-c.quit:
			gameServerInstance.Logout(c, "QUIT channel")
			return // terminate the client
		}
	}
}

func (c *Client) String() string {
	return fmt.Sprintf("(WebClient conn: %v, Player %s)", c.conn != nil, c.Player)
}
