package client

import (
	"github.com/trasa/watchmud/player"
)

//noinspection GoNameStartsWithPackageName
type Client interface {
	Send(innerMessage interface{}) error
	SetPlayer(player player.Player)
	GetPlayer() player.Player
	Close()
}
