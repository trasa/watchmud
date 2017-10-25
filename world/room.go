package world

import (
	"fmt"
	"github.com/trasa/watchmud/direction"
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/player"
	"github.com/trasa/watchmud/thing"
	"log"
	"math/rand"
	"strings"
	"time"
)

type Room struct {
	Id          string
	Name        string
	Description string
	Zone        *Zone
	PlayerList  *player.List // map of players by name
	Inventory   thing.Map
	Mobs        thing.Map
	North       *Room
	South       *Room
	East        *Room
	West        *Room
	Up          *Room
	Down        *Room
}

func NewRoom(zone *Zone, id string, name string, description string) *Room {
	return &Room{
		Id:          id,
		Name:        name,
		Description: description,
		Zone:        zone,
		PlayerList:  player.NewList(),
		Inventory:   make(thing.Map),
		Mobs:        make(thing.Map),
	}
}

func (r Room) String() string {
	return fmt.Sprintf("(Room %s: '%s')", r.Id, r.Name)
}

// Player leaves a room. Tells other room residents about it.
func (r *Room) PlayerLeaves(p player.Player, dir direction.Direction) {
	r.PlayerList.Remove(p)
	r.Send(message.LeaveRoomNotification{
		Response:  message.NewSuccessfulResponse("leave_room"),
		Name:      p.GetName(),
		Direction: dir,
	})
}

func (r *Room) MobileLeaves(mob *mobile.Instance, dir direction.Direction) {
	r.Mobs.Remove(mob)
	r.Send(message.LeaveRoomNotification{
		Response:  message.NewSuccessfulResponse("leave_room"),
		Name:      mob.Definition.Name, // TODO figure out name here...
		Direction: dir,
	})
}

// Add a player to a room. Don't send notifications.
func (r *Room) AddPlayer(p player.Player) {
	r.PlayerList.Add(p)
}

// Player enters a room. Tells other room residents about it.
func (r *Room) PlayerEnters(p player.Player) {
	r.Send(message.EnterRoomNotification{
		Response: message.NewSuccessfulResponse("enter_room"),
		Name:     p.GetName(),
	})
	r.AddPlayer(p)
}

func (r *Room) MobileEnters(mob *mobile.Instance) {
	r.Send(message.EnterRoomNotification{
		Response: message.NewSuccessfulResponse("enter_room"),
		Name:     mob.Definition.Name,
	})
	r.AddMobile(mob)
}

func (r *Room) AddMobile(inst *mobile.Instance) error {
	return r.Mobs.Add(inst)
}

// Send to every player in the room.
func (r *Room) Send(msg interface{}) { // TODO err
	r.PlayerList.Iter(func(p player.Player) {
		p.Send(msg)
	})
}

// Send to all players in a room, except for one of them.
func (r *Room) SendExcept(exception player.Player, msg interface{}) { // TODO err
	r.PlayerList.Iter(func(p player.Player) {
		if exception != p {
			p.Send(msg)
		}
	})
}

// Get all the valid exits from this room.
// TODO: exits can be locked and/or closed
func (r *Room) GetExitString() string {
	exits := []string{}
	r.forEachExit(exits, func(dir direction.Direction, context interface{}) {
		s, err := direction.DirectionToAbbreviation(dir)
		if err == nil {
			exits = append(exits, s)
		}
	})
	return strings.Join(exits, "")
}

func (r *Room) forEachExit(context interface{}, foreach func(dir direction.Direction, context interface{})) interface{} {
	if r.HasExit(direction.NORTH) {
		foreach(direction.NORTH, context)
	}
	if r.HasExit(direction.EAST) {
		foreach(direction.EAST, context)
	}
	if r.HasExit(direction.SOUTH) {
		foreach(direction.SOUTH, context)
	}
	if r.HasExit(direction.WEST) {
		foreach(direction.WEST, context)
	}
	if r.HasExit(direction.UP) {
		foreach(direction.UP, context)
	}
	if r.HasExit(direction.DOWN) {
		foreach(direction.DOWN, context)
	}
	return context
}

// Maps direction.Direction  to Name ("North Tower", "Garage")
type ExitInfo map[direction.Direction]string

// Return a map of info about the rooms around:
// directions that can be traveled and the short description of
// the room
// TODO some rooms can't be seen into, doors that are locked
// or closed, etc etc...
func (r *Room) GetExitInfo(limitToZone bool) ExitInfo {
	exits := make(ExitInfo)
	r.forEachExit(nil, func(dir direction.Direction, _ interface{}) {
		// TODO some rooms can't be seen into, etc ...
		if !limitToZone || r.Zone == r.Get(dir).Zone {
			exits[dir] = r.Get(dir).Name
		}
	})
	return exits
}

// Return a map of info about the rooms around this one:
// mapping the room IDs to the direction to take to get there.
func (r *Room) GetExitsByRoomId() map[string]direction.Direction {
	exits := make(map[string]direction.Direction)
	r.forEachExit(nil, func(dir direction.Direction, _ interface{}) {
		exits[r.Get(dir).Id] = dir
	})
	return exits
}

// Is there a valid exit in this direction in this room?
// TODO what about exits that are locked or closed?
func (r *Room) HasExit(dir direction.Direction) bool {
	switch dir {
	case direction.NORTH:
		return r.North != nil
	case direction.EAST:
		return r.East != nil
	case direction.SOUTH:
		return r.South != nil
	case direction.WEST:
		return r.West != nil
	case direction.UP:
		return r.Up != nil
	case direction.DOWN:
		return r.Down != nil
	default:
		log.Printf("room.HasExit: asked for unknown direction '%s'", dir)
		return false
	}
}

// Get the exit for this direction. Will return nil if there
// isn't a valid exit that way.
// TODO what about exits that are locked or closed?
func (r *Room) Get(dir direction.Direction) *Room {
	switch dir {
	case direction.NORTH:
		return r.North
	case direction.EAST:
		return r.East
	case direction.SOUTH:
		return r.South
	case direction.WEST:
		return r.West
	case direction.UP:
		return r.Up
	case direction.DOWN:
		return r.Down
	default:
		return nil
	}
}

func (r *Room) Set(dir direction.Direction, destRoom *Room) {
	switch dir {
	case direction.NORTH:
		r.North = destRoom
	case direction.EAST:
		r.East = destRoom
	case direction.SOUTH:
		r.South = destRoom
	case direction.WEST:
		r.West = destRoom
	case direction.UP:
		r.Up = destRoom
	case direction.DOWN:
		r.Down = destRoom
	default:
		log.Printf("Unknown direction to set: %s", dir)
	}
}

// Of the directions available for travel (could be locked, closed...)
// pick one of them. If there aren't any, return none.
func (r *Room) PickRandomDirection(limitToZone bool) direction.Direction {
	exits := r.GetExitInfo(limitToZone)
	if len(exits) == 0 {
		return direction.NONE
	} else {
		desired := rand.New(rand.NewSource(time.Now().Unix())).Int31n(int32(len(exits)))
		// iterate to the ith member of exits
		i := int32(0)
		for dir := range exits {
			if i == desired {
				return dir
			}
			i++
		}
		// inconceivable!
		log.Printf("Room.PickRandomDirection: Bizzare RandomDirection picked. len=%d, desired=%d", len(exits), desired)
		return direction.NONE
	}
}

// Describe this room.
func (r *Room) CreateRoomDescription(exclude player.Player) message.RoomDescription {
	desc := message.RoomDescription{
		Name:        r.Name,
		Description: r.Description,
		Exits:       r.GetExitString(),
	}
	// Note: the thread-safe iteration isn't necessary because only
	// one message is processed at a time (our server isn't actually
	// multithreaded...)
	r.PlayerList.Iter(func(p player.Player) {
		if p != exclude {
			desc.Players = append(desc.Players, p.GetName())
		}
	})
	for _, o := range r.Inventory {
		desc.Objects = append(desc.Objects, o.(*object.Instance).Definition.DescriptionOnGround)
	}
	for _, o := range r.Mobs {
		desc.Mobs = append(desc.Mobs, o.(*mobile.Instance).Definition.DescriptionInRoom)
	}

	return desc
}

func (r *Room) AddInventory(inst *object.Instance) error {
	return r.Inventory.Add(inst)
}
