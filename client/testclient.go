package client

import (
	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud/player"
)

type TestClient struct {
	Player player.Player
	tosend []interface{}
	open   bool
}

func NewTestClient(p player.Player) *TestClient {
	return &TestClient{
		Player: p,
	}
}

func (c *TestClient) Send(msg interface{}) error {
	log.Info().Msgf("sending fake! %s p is %s", msg, c.Player.GetName())
	c.tosend = append(c.tosend, msg)
	return nil
}

func (c *TestClient) GetPlayer() player.Player {
	return c.Player
}

func (c *TestClient) SetPlayer(p player.Player) {
	c.Player = p
}

func (c *TestClient) Close() {
	c.open = false
}
func (c *TestClient) GetSentResponse(i int) interface{} {
	return c.tosend[i]
}
