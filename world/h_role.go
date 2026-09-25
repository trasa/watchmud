package world

import (
	"fmt"

	"github.com/watchmud/watchmud/command"
	"github.com/watchmud/watchmud/event"
	"github.com/watchmud/watchmud/gameserver"
)

// handleRole shows the player what their equipment adds up to, and what it
// would take to add up to something else. There is no command to change a
// role -- wearing different gear is the command.
func (w *World) handleRole(msg *gameserver.HandlerParameter, cmd command.Role) {
	p := msg.Player

	// The reasons and the totals come off the same list, so a player can add
	// up what they are shown and get the number beside it. In slot order,
	// because their own gear should not reshuffle between asks.
	sources := make(map[string][]string)
	for _, c := range p.Equipment().RoleContributions() {
		sources[c.RoleId] = append(sources[c.RoleId],
			fmt.Sprintf("%s %d", c.Instance.Definition.ShortDescription, c.Weight))
	}

	weights := p.RoleWeights()
	standings := make([]event.RoleStanding, 0, len(w.content.Catalog.RoleList()))
	for _, r := range w.content.Catalog.RoleList() {
		standings = append(standings, event.RoleStanding{
			Name:    r.Name,
			Total:   weights[r.Id],
			Sources: sources[r.Id],
		})
	}

	e := event.Role{Standings: standings}
	if current := w.content.Catalog.RoleFor(weights); current != nil {
		e.Current = current.Name
		e.Description = current.Description
	}
	p.Send(e)
}
