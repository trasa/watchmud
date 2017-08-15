package world

import (
	"fmt"
	"github.com/trasa/watchmud/direction"
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/player"
	"log"
	"strings"
)

type Room struct {
	Id          string
	Name        string
	Description string
	Zone        *Zone        `json:"-"`
	PlayerList  *player.List `json:"-"` // map of players by name
	North       *Room        `json:"-"`
	South       *Room        `json:"-"`
	East        *Room        `json:"-"`
	West        *Room        `json:"-"`
	Up          *Room        `json:"-"`
	Down        *Room        `json:"-"`
}

func NewRoom(zone *Zone, id string, name string, description string) *Room {
	return &Room{
		Id:          id,
		Name:        name,
		Description: description,
		Zone:        zone,
		PlayerList:  player.NewList(),
	}
}

func (r Room) String() string {
	return fmt.Sprintf("(Room %s: '%s')", r.Id, r.Name)
}

// Player leaves a room. Tells other room residents about it.
func (r *Room) Leave(p player.Player, dir direction.Direction) {
	r.PlayerList.Remove(p)
	r.Send(message.LeaveRoomNotification{
		Notification: message.Notification{MessageType: "leave_room"},
		PlayerName:   p.GetName(),
		Direction:    dir,
	})
}

// Add a player to a room. Don't send notifications.
func (r *Room) Add(p player.Player) {
	r.PlayerList.Add(p)
}

// Player enters a room. Tells other room residents about it.
func (r *Room) Enter(p player.Player) {
	r.Send(message.EnterRoomNotification{
		Notification: message.Notification{MessageType: "enter_room"},
		PlayerName:   p.GetName(),
	})
	r.Add(p)
}

// Send to every player in the room.
func (r *Room) Send(msg interface{}) { // TODO err
	r.PlayerList.Iter(func(p player.Player) {
		p.Send(msg)
	})
}

// Get all the valid exits from this room.
// TODO: exits can be locked and/or closed
func (r *Room) GetExitString() string {
	exits := []string{}
	if r.HasExit(direction.NORTH) {
		exits = append(exits, "n")
	}
	if r.HasExit(direction.EAST) {
		exits = append(exits, "e")
	}
	if r.HasExit(direction.SOUTH) {
		exits = append(exits, "s")
	}
	if r.HasExit(direction.WEST) {
		exits = append(exits, "w")
	}
	if r.HasExit(direction.UP) {
		exits = append(exits, "u")
	}
	if r.HasExit(direction.DOWN) {
		exits = append(exits, "d")
	}
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

// Return a map of info about the rooms around:
// directions that can be traveled and the short description of
// the room
// TODO some rooms can't be seen into, doors that are locked
// or closed, etc etc...
func (r *Room) GetExitInfo() map[string]string {
	exits := make(map[string]string)
	r.forEachExit(exits, func(dir direction.Direction, context interface{}) {
		s, e := direction.DirectionToString(dir)
		if e == nil {
			// TODO some rooms can't be seen into, etc ...
			context.(map[string]string)[s] = r.Get(dir).Name
		} else {
			log.Printf("Couldn't DirectionToString: %s, %s", dir, e)
		}
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

// Describe this room.
func (r *Room) CreateRoomDescription() message.RoomDescription {
	return message.RoomDescription{
		Name:        r.Name,
		Description: r.Description,
		Exits:       r.GetExitString(),
		// TODO other occupants or objects in the room
	}
}
