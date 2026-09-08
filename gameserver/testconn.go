package gameserver

import "github.com/trasa/watchmud/player"

type TestConn struct {
	player *player.Player
	open   bool
}

func NewTestConn(p *player.Player) *TestConn {
	return &TestConn{
		player: p,
		open:   true,
	}
}

func (c *TestConn) Send(msg interface{}) error {
	return c.player.Send(msg)
}

func (c *TestConn) Player() *player.Player {
	return c.player
}

func (c *TestConn) SetPlayer(p *player.Player) {
	c.player = p
}

func (c *TestConn) Close() {
	c.open = false
}
