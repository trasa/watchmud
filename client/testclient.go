package client

import (
	"github.com/trasa/watchmud/player"
)

type TestClient struct {
	player *player.Player
	open   bool
}

func NewTestClient(p *player.Player) *TestClient {
	return &TestClient{
		player: p,
	}
}

func (c *TestClient) Send(msg interface{}) error {
	return c.player.Send(msg)
}

func (c *TestClient) Player() *player.Player {
	return c.player
}

func (c *TestClient) Close() {
	c.open = false
}
