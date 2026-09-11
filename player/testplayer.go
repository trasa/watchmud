package player

import (
	"uuid"

	"github.com/trasa/watchmud/rules"
)

// NewTestPlayer that tracks messages
func NewTestPlayer(id uuid.UUID, name string, out Sender) *Player {
	if out == nil {
		out = &Recorder{}
	}
	return New(id,
		name,
		out,
		&rules.Lineage{},
		&rules.Class{},
		rules.Abilities{})
}
