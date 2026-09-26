package spaces

import (
	"github.com/rs/zerolog/log"
	"github.com/watchmud/watchmud/mobile"
	"github.com/watchmud/watchmud/player"
	"github.com/watchmud/watchmud/rules"
)

// Occupancy is who is in which room, both directions, and the only code
// that writes a room's player and mob lists. The room's lists answer
// "who is here" in a stable order; the maps here are the index for
// "where is this one". Two copies are fine, two writers were the bug --
// see ROADMAP "Dual location bookkeeping".
type Occupancy struct {
	playerRoom map[*player.Player]*Room
	mobileRoom map[*mobile.Instance]*Room
}

func NewOccupancy() *Occupancy {
	return &Occupancy{
		playerRoom: make(map[*player.Player]*Room),
		mobileRoom: make(map[*mobile.Instance]*Room),
	}
}

// RoomOfPlayer is the room p is in, nil if they are nowhere.
func (o *Occupancy) RoomOfPlayer(p *player.Player) *Room {
	return o.playerRoom[p]
}

// PlacePlayer puts p in r without telling anyone: login and creation.
func (o *Occupancy) PlacePlayer(p *player.Player, r *Room) {
	if cur := o.playerRoom[p]; cur != nil {
		log.Error().Str("player", p.Name()).Str("room", cur.Location().String()).
			Msg("PlacePlayer: already placed")
		return
	}
	r.AddPlayer(p)
	o.playerRoom[p] = r
}

// MovePlayer takes p out of wherever they are and into dest, telling both
// rooms. A player who is nowhere just arrives.
func (o *Occupancy) MovePlayer(p *player.Player, dir rules.Direction, dest *Room) {
	if src := o.playerRoom[p]; src != nil {
		src.PlayerLeaves(p, dir)
	}
	dest.PlayerEnters(p)
	o.playerRoom[p] = dest
}

// RemovePlayer takes p out of the world's rooms without telling anyone.
func (o *Occupancy) RemovePlayer(p *player.Player) {
	if r := o.playerRoom[p]; r != nil {
		r.RemovePlayer(p)
	}
	delete(o.playerRoom, p)
}

func (o *Occupancy) RoomOfMobile(mob *mobile.Instance) *Room {
	return o.mobileRoom[mob]
}

// PlaceMobile puts mob in r without telling anyone: zone resets and `load`.
func (o *Occupancy) PlaceMobile(mob *mobile.Instance, r *Room) {
	if cur := o.mobileRoom[mob]; cur != nil {
		log.Error().Str("mob", mob.Name()).Str("room", cur.Location().String()).
			Msg("PlaceMobile: already placed")
		return
	}
	if err := r.AddMobile(mob); err != nil {
		log.Error().Err(err).Str("room", r.Location().String()).Msg("PlaceMobile")
		return
	}
	o.mobileRoom[mob] = r
}

func (o *Occupancy) MoveMobile(mob *mobile.Instance, dir rules.Direction, dest *Room) {
	if src := o.mobileRoom[mob]; src != nil {
		src.MobileLeaves(mob, dir)
	}
	dest.MobileEnters(mob)
	o.mobileRoom[mob] = dest
}

func (o *Occupancy) RemoveMobile(mob *mobile.Instance) {
	if r := o.mobileRoom[mob]; r != nil {
		if err := r.RemoveMobile(mob); err != nil {
			log.Error().Err(err).Str("room", r.Location().String()).Msg("RemoveMobile")
		}
	}
	delete(o.mobileRoom, mob)
}

func (o *Occupancy) Mobiles() []*mobile.Instance {
	mobs := make([]*mobile.Instance, 0, len(o.mobileRoom))
	for m := range o.mobileRoom {
		mobs = append(mobs, m)
	}
	return mobs
}

func (o *Occupancy) MobileCount(defId string) int {
	count := 0
	for mob := range o.mobileRoom {
		if mob.Definition.Id == defId {
			count++
		}
	}
	return count
}
