package server

import "log"

type TestClient struct {
	Player *Player
	tosend []interface{}
	open   bool
}

func (c *TestClient) Send(msg interface{}) {
	log.Printf("sending fake! %s p is %s", msg, c.Player.Name)
	c.tosend = append(c.tosend, msg)
}

func (c *TestClient) GetPlayer() *Player {
	return c.Player
}

func (c *TestClient) SetPlayer(p *Player) {
	c.Player = p
}

func (c *TestClient) Close() {
	c.open = false
}
