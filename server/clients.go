package server

import (
	"log"
	"sync"
)

// collection of clients
var clients = newClients()

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
