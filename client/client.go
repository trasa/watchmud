package client

import "github.com/trasa/watchmud/player"

type Client interface {
	Send(message interface{}) // todo return err
	SetPlayer(player player.Player)
	GetPlayer() player.Player
	Close()
}
