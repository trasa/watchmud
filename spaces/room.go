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
	mobs        *RoomMobs
	directions  map[direction.Direction]*Room
	flags       map[string]bool
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
		mobs:        NewRoomMobs(),
		directions:  make(map[direction.Direction]*Room),
		flags:       make(map[string]bool),
	}
}

// Build a strip down version of a Room, for testing
func NewTestRoom(name string) *Room {
	return NewRoom(nil, name, name, "")
}

func (r Room) String() string {
	return fmt.Sprintf("(Room %s: '%s')", r.Id, r.Name)
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

func (r *Room) HasFlag(flag string) bool {
	return r.flags[flag]
}

// Player leaves a room. Tells other room residents about it.
func (r *Room) PlayerLeaves(p player.Player, dir direction.Direction) {
	r.playerList.Remove(p)
	r.Send(message.LeaveRoomNotification{
		Success:    true,
		ResultCode: "OK",
		Name:       p.GetName(),
		Direction:  int32(dir),
	})
}

func (r *Room) MobileLeaves(mob *mobile.Instance, dir direction.Direction) {
	r.mobs.Remove(mob)
	r.Send(message.LeaveRoomNotification{
		Success:    true,
		ResultCode: "OK",
		Name:       mob.Definition.Name, // TODO figure out name here...
		Direction:  int32(dir),
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
		Success:    true,
		ResultCode: "OK",
		Name:       p.GetName(),
	})
	r.AddPlayer(p)
}

func (r *Room) MobileEnters(mob *mobile.Instance) {
	r.Send(message.EnterRoomNotification{
		Success:    true,
		ResultCode: "OK",
		Name:       mob.Definition.Name,
	})
	r.AddMobile(mob)
}

func (r *Room) AddMobile(inst *mobile.Instance) error {
	return r.mobs.Add(inst)
}

func (r *Room) RemoveMobile(inst *mobile.Instance) error {
	return r.mobs.Remove(inst)
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
func (r *Room) CreateRoomDescription(exclude player.Player) *message.RoomDescription {
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
	for _, mob := range r.mobs.GetAll() {
		desc.Mobs = append(desc.Mobs, mob.Definition.DescriptionInRoom)
	}
	return &desc
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

// Attempt to find an inventory object in this room which matches the
// terms given. Search the object name and aliases.
func (r *Room) FindInventory(findMode message.FindMode, index string, target string) (inst *object.Instance, exists bool) {
	inst, exists = r.inventory.Find(findMode, index, target)
	return
}

func (r *Room) FindMobile(target string) (mob *mobile.Instance, exists bool) {
	return r.mobs.Find(target)
}
