package player

import (
	"uuid"

	"github.com/trasa/watchmud/rules"
)

// NewTestPlayer that tracks messages. The test catalog's armor table is
// empty, so equipment is worth exactly what it declares by hand -- a test that
// wants armor to argue for a role builds its own catalog.
func NewTestPlayer(id uuid.UUID, name string, out Sender) *Player {
	if out == nil {
		out = &Recorder{}
	}
	cat, err := rules.NewTestCatalog()
	if err != nil {
		panic(err) // the test catalog is a literal; it cannot fail to index
	}
	return New(
		id,
		name,
		"passwordHash",
		out,
		&rules.Lineage{Id: "human", Name: "Human"},
		cat,
	)
}
