package world

import (
	"github.com/trasa/watchmud/response"
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

func (p *TestPlayer) Send(msg interface{}) {
	log.Printf("sending for player %s", msg)
	p.sent = append(p.sent, msg)
}

func (p *TestPlayer) GetName() string {
	return p.name
}

func (p *TestPlayer) SentMessageCount() int {
	return len(p.sent)
}

func (p *TestPlayer) GetSentResponse(i int) response.Response {
	return p.sent[i].(response.Response)
}
