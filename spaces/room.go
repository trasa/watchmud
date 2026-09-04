package spaces

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud-message/direction"
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

func (r *Room) GetFlags() (result []string) {
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
	r.Send(message.LeaveRoomNotification{
		Success:    true,
		ResultCode: "OK",
		Name:       p.Name,
		Direction:  int32(dir),
	})
}

func (r *Room) MobileLeaves(mob *mobile.Instance, dir direction.Direction) {
	// TODO unhandled error
	r.mobs.Remove(mob)
	r.Send(message.LeaveRoomNotification{
		Success:    true,
		ResultCode: "OK",
		Name:       mob.Definition.Name, // TODO figure out name here...
		Direction:  int32(dir),
	})
}

// AddPlayer to the Room, without sending notifications
func (r *Room) AddPlayer(p *player.Player) {
	r.playerList.Add(p)
}

// RemovePlayer from the Room, without sending notifications
func (r *Room) RemovePlayer(p *player.Player) {
	r.playerList.Remove(p)
}

func (r *Room) GetPlayers() []*player.Player {
	return r.playerList.GetAll()
}

// PlayerEnters a room, telling other room entities about it.
func (r *Room) PlayerEnters(p *player.Player) {
	// TODO how does this make sense next to the other Add, etc funcs?
	r.Send(message.EnterRoomNotification{
		Success:    true,
		ResultCode: "OK",
		Name:       p.Name,
	})
	r.AddPlayer(p)
}

// MobileEnters a room, telling other room entities about it.
func (r *Room) MobileEnters(mob *mobile.Instance) {
	// TODO how does this make sense next to the other Add, etc funcs?
	r.Send(message.EnterRoomNotification{
		Success:    true,
		ResultCode: "OK",
		Name:       mob.Definition.Name,
	})
	// TODO error handling
	r.AddMobile(mob)
}

func (r *Room) AddMobile(inst *mobile.Instance) error {
	return r.mobs.Add(inst)
}

func (r *Room) RemoveMobile(inst *mobile.Instance) error {
	return r.mobs.Remove(inst)
}

func (r *Room) GetMobs() []*mobile.Instance {
	return r.mobs.GetAll()
}

// Send to every player in the room.
func (r *Room) Send(msg interface{}) {
	r.playerList.Iter(func(p *player.Player) {
		// TODO error handling
		p.Send(msg)
	})
}

// SendExcept to one player
func (r *Room) SendExcept(exception *player.Player, msg interface{}) {
	r.playerList.Iter(func(p *player.Player) {
		if exception != p {
			// TODO error handling
			p.Send(msg)
		}
	})
}

// Notify everything in a room about something
func (r *Room) Notify(msg interface{}) {
	// TODO error handling
	for _, m := range r.mobs.GetAll() {
		m.Send(msg)
	}
	r.Send(msg)
}

// CreateRoomDescription describes the room, except for one player
func (r *Room) CreateRoomDescription(exclude *player.Player) *message.RoomDescription {
	// TODO rework this
	desc := message.RoomDescription{
		Name:        r.Name,
		Description: r.Description,
		Exits:       r.GetExitString(),
	}
	// Note: the thread-safe iteration isn't necessary because only
	// one message is processed at a time (our server isn't actually
	// multithreaded...)
	r.playerList.Iter(func(p *player.Player) {
		if p != exclude {
			desc.Players = append(desc.Players, p.Name)
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

// GetInventoryByName finds the object.Instance matching by name
func (r *Room) GetInventoryByName(name string) (inst *object.Instance, exists bool) {
	inst, exists = r.inventory.GetByName(name)
	return
}

func (r *Room) GetAllInventory() []*object.Instance {
	return r.inventory.GetAll()
}

// FindInventory in the room matching terms given. Searches the object names and aliases.
func (r *Room) FindInventory(findMode message.FindMode, index string, target string) (inst *object.Instance, exists bool) {
	inst, exists = r.inventory.Find(findMode, index, target)
	return
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
