package spaces

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud/direction"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/player"
)

type Room struct {
	Id          string
	Name        string
	Description string
	Zone        *Zone
	playerList  *player.List // map of players by name
	Inventory   *RoomInventory
	mobs        *RoomMobs
	directions  map[direction.Direction]*Room
	flags       map[string]bool
}

// NewRoom creates a new room in this zone.
// It does NOT alter the Zone to include the new room, however.
func NewRoom(zone *Zone, id string, name string, description string) *Room {
	return &Room{
		Id:          id,
		Name:        name,
		Description: description,
		Zone:        zone,
		playerList:  player.NewList(),
		Inventory:   NewRoomInventory(),
		mobs:        NewRoomMobs(),
		directions:  make(map[direction.Direction]*Room),
		flags:       make(map[string]bool),
	}
}

// NewTestRoom creates a simpler room for tests
func NewTestRoom(name string) *Room {
	return NewRoom(nil, name, name, "")
}

func (r *Room) String() string {
	return fmt.Sprintf("(Room %s-%s: '%s')", r.Zone.Id, r.Id, r.Name)
}

func (r *Room) Location() player.Location {
	return player.NewLocation(r.Zone.Id, r.Id)
}

func (r *Room) SetFlags(flags []string) {
	if flags != nil {
		for _, s := range flags {
			r.SetFlag(s)
		}
	}
}

func (r *Room) SetFlag(flag string) {
	r.flags[flag] = true
}

func (r *Room) Flag(flag string) bool {
	return r.flags[flag]
}

func (r *Room) Flags() (result []string) {
	for k, v := range r.flags {
		if v {
			result = append(result, k)
		}
	}
	return
}

// PlayerLeaves a room. Tells other room residents about it.
func (r *Room) PlayerLeaves(p *player.Player, dir direction.Direction) {
	r.playerList.Remove(p)
	r.Send(event.Left{Who: p.Name(), Direction: dir})
}

func (r *Room) MobileLeaves(mob *mobile.Instance, dir direction.Direction) {
	if err := r.mobs.Remove(mob); err != nil {
		log.Error().Err(err).Str("room", r.Location().String()).Msg("mobileLeaves: failed to leave room")
		return
	}
	r.Send(event.Left{Who: mob.Name(), Direction: dir})
}

// AddPlayer to the Room, without sending notifications
func (r *Room) AddPlayer(p *player.Player) {
	r.playerList.Add(p)
}

// RemovePlayer from the Room, without sending notifications
func (r *Room) RemovePlayer(p *player.Player) {
	r.playerList.Remove(p)
}

func (r *Room) Players() []*player.Player {
	return r.playerList.GetAll()
}

func (r *Room) playersExcept(exclude *player.Player) []*player.Player {
	return r.playerList.GetExcept(exclude)
}

// PlayerEnters a room, telling other room entities about it.
func (r *Room) PlayerEnters(p *player.Player) {
	// TODO how does this make sense next to the other Add, etc funcs?
	r.Send(event.Entered{Who: p.Name()})
	r.AddPlayer(p)
}

// MobileEnters a room, telling other room entities about it.
func (r *Room) MobileEnters(mob *mobile.Instance) {
	if err := r.AddMobile(mob); err != nil {
		log.Error().Err(err).Str("room", r.Location().String()).Msg("MobileEnters: failed to add mobile")
		return
	}
	r.Send(event.Entered{Who: mob.Definition.Name})
}

func (r *Room) AddMobile(inst *mobile.Instance) error {
	return r.mobs.Add(inst)
}

func (r *Room) RemoveMobile(inst *mobile.Instance) error {
	return r.mobs.Remove(inst)
}

func (r *Room) Mobs() []*mobile.Instance {
	return r.mobs.GetAll()
}

// Send to every player in the room.
func (r *Room) Send(msg any) {
	r.playerList.Iter(func(p *player.Player) {
		p.Send(msg)
	})
}

// SendExcept to one player
func (r *Room) SendExcept(exception *player.Player, msg any) {
	r.playerList.Iter(func(p *player.Player) {
		if exception != p {
			p.Send(msg)
		}
	})
}

// Notify mobs and players in a room about something
func (r *Room) Notify(msg any) {
	// mobs
	for _, m := range r.mobs.GetAll() {
		m.Send(msg)
	}
	// players
	r.Send(msg)
}

// DescriptionExcept describes the room, except for one player, if provided
// The exits are stringified for telnet clients, not the abbreviations.
func (r *Room) DescriptionExcept(exclude *player.Player) event.RoomDescription {
	desc := event.RoomDescription{
		Name:        r.Name,
		Description: r.Description,
		Exits:       r.ExitString(),
	}

	for _, p := range r.playerList.GetExcept(exclude) {
		desc.Players = append(desc.Players, p.Name())
	}

	for _, o := range r.Inventory.GetAll() {
		desc.Objects = append(desc.Objects, o.Definition.DescriptionOnGround)
	}
	for _, mob := range r.mobs.GetAll() {
		desc.Mobs = append(desc.Mobs, mob.Definition.DescriptionInRoom)
	}
	return desc
}

func (r *Room) FindMobile(target string) (mob *mobile.Instance, exists bool) {
	return r.mobs.Find(target)
}

func (r *Room) FindPlayer(target string) (*player.Player, bool) {
	p := r.playerList.FindByName(target)
	if p != nil {
		return p, true
	}
	return nil, false
}

// ExitString returns all the valid exits from this room as a string.
func (r *Room) ExitString() string {
	// TODO: exits can be locked and/or closed, this doesn't handle that.
	var exits []direction.Direction
	for _, exit := range r.Exits(false) {
		exits = append(exits, exit.Direction)
	}
	return direction.Format(exits)
}

// HasExit determines if there is a valid exit in this direction.
func (r *Room) HasExit(dir direction.Direction) bool {
	// TODO what about exits that are locked or closed?
	_, ok := r.directions[dir]
	return ok
}

// DestinationRoom returns the room in this direction or nil if there isn't one.
func (r *Room) DestinationRoom(dir direction.Direction) (dest *Room) {
	// TODO what about exits that are locked or closed?
	return r.directions[dir]
}

// Connect this room to the destination room in this direction.
// Loader use only: room topology is immutable once content is loaded.
// See ROADMAP "Known Problems": room conflates definition and instance.
func (r *Room) Connect(dir direction.Direction, destRoom *Room) {
	r.directions[dir] = destRoom
}

// PickRandomDirection from here to travel, from directions that are available.
// If there aren't any, return direction.None.
func (r *Room) PickRandomDirection(limitToZone bool) direction.Direction {
	// TODO should this return a room and not a direction?
	exits := r.Exits(limitToZone)
	if len(exits) == 0 {
		return direction.None
	} else {
		desired := rand.New(rand.NewSource(time.Now().Unix())).Int31n(int32(len(exits)))
		// iterate to the ith member of exits
		i := int32(0)
		for _, re := range exits {
			if i == desired {
				return re.Direction
			}
			i++
		}
		// inconceivable!
		log.Warn().Msgf("Room.PickRandomDirection: Bizarre RandomDirection picked. len=%d, desired=%d", len(exits), desired)
		return direction.None
	}
}

// Exits returns the exits from this room.
// Uses the direction.Direction ordering.
// Does not take locks, doors, closures, etc. into account.
func (r *Room) Exits(limitToZone bool) []RoomExit {
	holder := roomExitHolder{}
	for dir, dest := range r.directions {
		if !limitToZone || r.Zone == dest.Zone {
			holder.dirs = append(holder.dirs, RoomExit{dir, dest})
		}
	}
	sort.Sort(holder)
	return holder.dirs
}
