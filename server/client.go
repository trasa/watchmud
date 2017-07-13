package server

import (
	"github.com/gorilla/websocket"
	"github.com/trasa/watchmud/world"
	"log"
	"sync"
)

// controls terminating all clients
var GlobalQuit = make(chan interface{})

// collection of clients
var clients = newClients()

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
	// websocket connection
	conn *websocket.Conn
	// buffered channel of outbound messages
	source chan interface{} // sends down to client
	quit   chan interface{} // used to terminate clients
	Player *world.Player
}

func newClient(c *websocket.Conn) *Client {
	return &Client{
		conn:   c,
		source: make(chan interface{}, 256),
		quit:   GlobalQuit,
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
			clients.remove(c)
			return
		}
		log.Printf("message body: %s", body)
		GameServerInstance.incomingMessageBuffer <- newIncomingMessage(c, body)
		//c.send <- []byte("abcdefg")
	}
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()
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

type Clients struct {
	sync.Mutex
	clients map[*Client]*Client
}

func newClients() *Clients {
	return &Clients{
		clients: make(map[*Client]*Client),
	}
}

func (cs *Clients) add(c *Client) {
	cs.Lock()
	defer cs.Unlock()
	cs.clients[c] = c
}

func (cs *Clients) remove(c *Client) {
	cs.Lock()
	defer cs.Unlock()
	delete(cs.clients, c)
}

func (cs *Clients) iter(routine func(*Client)) {
	cs.Lock()
	defer cs.Unlock()
	log.Printf("sending to %d clients", len(cs.clients))
	for c := range cs.clients {
		routine(c)
	}
}
