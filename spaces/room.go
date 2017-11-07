package spaces

import (
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/trasa/watchmud/direction"
	"github.com/trasa/watchmud/message"
	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/player"
)

type Room struct {
	Id          string
	Name        string
	Description string
	Zone        *Zone
	playerList  *player.List // map of players by name
	inventory   *RoomInventory
	mobs        map[*mobile.Instance]bool
	directions  map[direction.Direction]*Room
}

// Create a new Room reference
func NewRoom(zone *Zone, id string, name string, description string) *Room {
	return &Room{
		Id:          id,
		Name:        name,
		Description: description,
		Zone:        zone,
		playerList:  player.NewList(),
		inventory:   NewRoomInventory(),
		mobs:        make(map[*mobile.Instance]bool),
		directions:  make(map[direction.Direction]*Room),
	}
}

// Build a strip down version of a Room, for testing
func NewTestRoom(name string) *Room {
	return NewRoom(nil, name, name, "")
}

func (r Room) String() string {
	return fmt.Sprintf("(Room %s: '%s')", r.Id, r.Name)
}

// Player leaves a room. Tells other room residents about it.
func (r *Room) PlayerLeaves(p player.Player, dir direction.Direction) {
	r.playerList.Remove(p)
	r.Send(message.LeaveRoomNotification{
		Response:  message.NewSuccessfulResponse("leave_room"),
		Name:      p.GetName(),
		Direction: dir,
	})
}

func (r *Room) MobileLeaves(mob *mobile.Instance, dir direction.Direction) {
	//r.Mobs.Remove(mob)
	r.mobs[mob] = false
	r.Send(message.LeaveRoomNotification{
		Response:  message.NewSuccessfulResponse("leave_room"),
		Name:      mob.Definition.Name, // TODO figure out name here...
		Direction: dir,
	})
}

// Add a player to a room. Don't send notifications.
func (r *Room) AddPlayer(p player.Player) {
	r.playerList.Add(p)
}

func (r *Room) RemovePlayer(p player.Player) {
	r.playerList.Remove(p)
}

func (r *Room) GetPlayers() []player.Player {
	return r.playerList.GetAll()
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
	//return r.Mobs.Add(inst)
	r.mobs[inst] = true
	return nil
}

func (r *Room) RemoveMobile(inst *mobile.Instance) {
	r.mobs[inst] = false
}

// Send to every player in the room.
func (r *Room) Send(msg interface{}) {
	r.playerList.Iter(func(p player.Player) {
		p.Send(msg)
	})
}

// Send to all players in a room, except for one of them.
func (r *Room) SendExcept(exception player.Player, msg interface{}) {
	r.playerList.Iter(func(p player.Player) {
		if exception != p {
			p.Send(msg)
		}
	})
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
	r.playerList.Iter(func(p player.Player) {
		if p != exclude {
			desc.Players = append(desc.Players, p.GetName())
		}
	})
	for _, o := range r.inventory.GetAll() {
		desc.Objects = append(desc.Objects, o.Definition.DescriptionOnGround)
	}
	for mob, inroom := range r.mobs {
		if inroom {
			desc.Mobs = append(desc.Mobs, mob.Definition.DescriptionInRoom)
		}
	}
	return desc
}

func (r *Room) AddInventory(inst *object.Instance) error {
	return r.inventory.Add(inst)
}

func (r *Room) RemoveInventory(inst *object.Instance) error {
	return r.inventory.Remove(inst)
}

func (r *Room) GetInventoryByInstanceId(instanceId uuid.UUID) (inst *object.Instance, exists bool) {
	inst, exists = r.inventory.GetByInstanceId(instanceId)
	return
}

// Find an object.Instance that matches this name
func (r *Room) GetInventoryByName(name string) (inst *object.Instance, exists bool) {
	inst, exists = r.inventory.GetByName(name)
	return
}

func (r *Room) GetAllInventory() []*object.Instance {
	return r.inventory.GetAll()
}
