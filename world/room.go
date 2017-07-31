package world

import (
	"fmt"
	"github.com/trasa/watchmud/direction"
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

func (r *Room) RemovePlayer(p player.Player) {
	r.PlayerList.Remove(p)
	// tell players in room that this player has left
	/*
		for _, p := range r.Players {
			p.OnEvent(ExitRoomEvent{
				PlayerId: player.Id,
				ZoneId:   r.Zone.Id,
				RoomId:   r.Id,
			})
		}
	*/
}

func (r *Room) AddPlayer(p player.Player) {
	/*
		for _, p := range r.Players {
			p.OnEvent(EnterRoomEvent{
				PlayerId: player.Id,
				ZoneId:   r.Zone.Id,
				RoomId:   r.Id,
			})
		}
	*/
	r.PlayerList.Add(p)
}

func (r *Room) GetExits() string {
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

type ExitRoomEvent struct {
	PlayerId string
	ZoneId   string
	RoomId   string
}

type EnterRoomEvent struct {
	PlayerId string
	ZoneId   string
	RoomId   string
}
