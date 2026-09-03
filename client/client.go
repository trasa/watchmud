package client

import (
	"github.com/trasa/watchmud/player"
)

// noinspection GoNameStartsWithPackageName
type Client interface {
	Send(innerMessage any) error
	SetPlayer(player player.Player)
	GetPlayer() player.Player
	Close()
}
