package event

import (
	"github.com/trasa/watchmud/direction"
	"github.com/trasa/watchmud/slot"
)

// ---- session ---------------------------------------------------------------
//
// These four are read by the connection's login conversation
// (telnet/conn.go), not by the renderer. They are four types rather than two
// with a Success field because a type switch there stays honest in a way a
// bool doesn't.

type LoggedIn struct {
	Name string
}

type LoginFailed struct {
	Reason ResultCode
}

type PlayerCreated struct {
	Name string
}

type CreateFailed struct {
	Reason ResultCode
}

// LoggedOut reaches the room the player just left, so in practice the player
// themselves never sees it -- they are removed from the room first, and their
// own goodbye comes from the connection.
type LoggedOut struct {
	Actor string
}

type Pong struct {
	Target string
}

// ---- rooms -----------------------------------------------------------------

// RoomDescription is what a player sees of a room. One event covers look,
// move and recall, which have always rendered identically.
type RoomDescription struct {
	Name        string
	Description string
	Exits       string
	Players     []string
	Objects     []string
	Mobs        []string
}

type Entered struct {
	Who string
}

type Left struct {
	Who       string
	Direction direction.Direction
}

type Exits struct {
	Exits []Exit
}

type Exit struct {
	Direction direction.Direction
	RoomName  string
}

// ---- objects ---------------------------------------------------------------
//
// Dropped and Got are each one event for both audiences: the player who acted
// and everyone else in the room. The renderer compares Actor to the name of
// the player it is rendering for, the same way it already does for Struck.

type Dropped struct {
	Actor string
	Item  string
}

type Got struct {
	Actor string
	Item  string
}

// Equipped and Worn have no bystander text today, so they carry no Actor.
type Equipped struct{}

type Worn struct{}

type Inventory struct {
	Items []InventoryItem
}

type InventoryItem struct {
	Id               string
	ShortDescription string
	Categories       []string
}

type Equipment struct {
	Items []EquippedItem
}

type EquippedItem struct {
	Slot             slot.Location
	Id               string
	ShortDescription string
}

// ---- talking ---------------------------------------------------------------

type Said struct {
	Speaker string
	Value   string
}

// Told goes to both parties: the renderer shows the sender an acknowledgement
// and the receiver the message.
type Told struct {
	From  string
	To    string
	Value string
}

type Shouted struct {
	Speaker string
	Value   string
}

// ---- the player -----------------------------------------------------------

type Who struct {
	Players []WhoEntry
}

type WhoEntry struct {
	PlayerName string
	ZoneName   string
	RoomName   string
}

type Stat struct {
	PlayerName    string
	Lineage       string
	Class         string
	CurrentHealth int
	MaxHealth     int
	ZoneId        string
	RoomId        string
	Strength      int
	Dexterity     int
	Constitution  int
	Intelligence  int
	Wisdom        int
	Charisma      int
}

// ---- combat ----------------------------------------------------------------

// Attacking acknowledges that a fight has started.
type Attacking struct {
	Target string
}

// Struck is one swing, seen by the whole room. The renderer picks the second
// person for whichever end of it is reading.
type Struck struct {
	Attacker string
	Target   string
	Hit      bool
	Damage   int
}

type Died struct {
	Target   string
	IsPlayer bool
}

type Restored struct {
	Target   string
	IsPlayer bool
}

// ---- builder commands ------------------------------------------------------

type Loaded struct{}

type RoomStatus struct {
	Id          string
	Name        string
	Description string
	ZoneId      string
	ZoneName    string
	Flags       []string
	Players     []RoomStatusPlayer
	Items       []RoomStatusItem
	Mobs        []RoomStatusMob
	Exits       []RoomStatusExit
}

type RoomStatusPlayer struct {
	Name          string
	CurrentHealth int
	MaxHealth     int
}

type RoomStatusItem struct {
	Id                  string
	DefinitionId        string
	ZoneId              string
	Name                string
	ShortDescription    string
	DescriptionOnGround string
	Aliases             []string
	Categories          []string
	Behaviors           []string
}

type RoomStatusMob struct {
	Id                string
	DefinitionId      string
	ZoneId            string
	Name              string
	ShortDescription  string
	DescriptionInRoom string
	Aliases           []string
	Flags             []string
	CurrentHealth     int
	MaxHealth         int
}

type RoomStatusExit struct {
	Direction direction.Direction
	RoomId    string
	ZoneId    string
	Flags     []string
}
