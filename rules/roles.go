package rules

import "fmt"

// A Role is what a character is doing in a fight: holding the line, keeping
// people alive, putting things down. It replaces the idea of a class.
//
// A role is never chosen and never stored. It is read off the equipment the
// character is wearing right now: every object definition can declare what it
// contributes to one or more roles (object.Definition.RoleWeights), the
// contributions of everything equipped are summed, and the highest total
// wins. Hang armor off every slot and you are a Tank; swap it for a censer
// and a prayer book and you are a Healer, in the time it takes to wear them.
//
// This is why there is no Role field on player.Player and no RoleId in
// player.Record. Storing it would let it disagree with the gear, and the gear
// is the truth.
type Role struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// indexRoles builds the lookup map, failing on duplicate or missing ids. The
// returned slice preserves the order roles were declared in roles.json, which
// is what breaks ties in RoleFor and what orders the `role` listing.
func indexRoles(all []*Role) (map[string]*Role, []*Role, error) {
	roles := make(map[string]*Role, len(all))
	order := make([]*Role, 0, len(all))
	for _, r := range all {
		if r.Id == "" {
			return nil, nil, fmt.Errorf("role %q: missing id", r.Name)
		}
		if prev, dup := roles[r.Id]; dup {
			return nil, nil, fmt.Errorf("duplicate role id %q (%s and %s)", r.Id, prev.Name, r.Name)
		}
		roles[r.Id] = r
		order = append(order, r)
	}
	return roles, order, nil
}
