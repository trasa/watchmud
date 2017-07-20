package world

import (
	"log"
)

type TestPlayer struct {
	name string
	sent []interface{}
}

// create a new test player that can track sent messages through 'sentmessages'
func NewTestPlayer(name string) *TestPlayer {
	p := &TestPlayer{
		name: name,
	}
	return p
}

func (this *TestPlayer) Send(msg interface{}) {
	log.Printf("sending for player %s", msg)
	this.sent = append(this.sent, msg)
}

func (this *TestPlayer) GetName() string {
	return this.name
}
