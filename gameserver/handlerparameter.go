package gameserver

import (
	"github.com/watchmud/watchmud/command"
	"github.com/watchmud/watchmud/event"
	"github.com/watchmud/watchmud/player"
)

// HandlerParameter is what a handler is handed: the connection, the player
// that connection had when the command arrived, and the command itself.
type HandlerParameter struct {
	Client  Conn
	Player  *player.Player
	Command command.Command
}

func NewHandlerParameter(c Conn, cmd command.Command) *HandlerParameter {
	return &HandlerParameter{
		Client:  c,
		Player:  c.Player(),
		Command: cmd,
	}
}

// Fail tells the player their command didn't work. The verb comes from the
// command itself, so a handler never has to repeat its own name.
func (h *HandlerParameter) Fail(code event.ResultCode) {
	if h.Player == nil {
		return
	}
	verb := ""
	if h.Command != nil {
		verb = h.Command.Verb()
	}
	h.Player.Send(event.Failed{Verb: verb, Code: code})
}
