package world

import (
	"fmt"
	"maps"
	"slices"

	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/gameserver"
)

// handleRole shows the player what their equipment adds up to, and what it
// would take to add up to something else. There is no command to change a
// role -- wearing different gear is the command.
func (w *World) handleRole(msg *gameserver.HandlerParameter, cmd command.Role) {
	p := msg.Player
	equipped := p.Slots().GetAll()

	// Slot order, because GetAll hands back a map and a player should not see
	// their own gear listed in a different order every time they ask.
	sources := make(map[string][]string)
	for _, loc := range slices.Sorted(maps.Keys(equipped)) {
		inst := equipped[loc]
		if inst == nil {
			continue // a slot whose item didn't survive a content edit
		}
		for roleId, weight := range inst.Definition.RoleWeights {
			if weight == 0 {
				continue
			}
			sources[roleId] = append(sources[roleId],
				fmt.Sprintf("%s %d", inst.Definition.ShortDescription, weight))
		}
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
