package gameserver

import "github.com/trasa/watchmud/player"

// Conn is a live connection to something that can play the game. It exists
// because a connection has no *player.Player until it logs in; once it does,
// everything talks through the player instead.
type Conn interface {
	player.Sender
	Player() *player.Player
	SetPlayer(p *player.Player)
	Close()
}
