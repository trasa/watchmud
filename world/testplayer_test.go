package world

import (
	"github.com/trasa/watchmud/client"
	"log"
)

type TestPlayer struct {
	name   string
	client client.Client
}

// create a new test player that can track sent messages through 'sentmessages'
func NewTestPlayer(name string) *TestPlayer {
	c := &TestClient{}
	p := &TestPlayer{
		name:   name,
		client: c,
	}
	c.Player = p
	return p
}

func (this *TestPlayer) Send(message interface{}) {
	log.Printf("sending for player %s", message)
	this.client.Send(message)
}

func (this *TestPlayer) GetName() string {
	return this.name
}
