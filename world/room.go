package world

import (
	"fmt"
	"github.com/trasa/watchmud/player"
	"strings"
)

type Room struct {
	Id          string
	Name        string
	Description string
	Zone        *Zone        `json:"-"`
	PlayerList  *player.List `json:"-"` // map of players by name
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
	if r.North != nil {
		exits = append(exits, "n")
	}
	if r.East != nil {
		exits = append(exits, "e")
	}
	if r.South != nil {
		exits = append(exits, "s")
	}
	if r.West != nil {
		exits = append(exits, "w")
	}
	if r.Up != nil {
		exits = append(exits, "u")
	}
	if r.Down != nil {
		exits = append(exits, "d")
	}
	return strings.Join(exits, "")
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
