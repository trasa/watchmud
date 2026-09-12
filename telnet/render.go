package telnet

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud/direction"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/player"
	"github.com/trasa/watchmud/rules"
)

// render turns anything sent to a connection into the text a telnet client
// sees. It is the single chokepoint between the game's event vocabulary and
// the wire, so keep game logic out of it.
//
// self is the name of the player this connection belongs to. Events that the
// whole room sees arrive here once per player, and the case compares Actor to
// self to pick between "Dropped." and "bob drops a knife." -- which is why
// there is no separate notification type for them.
func render(msg any, self string) string {
	switch m := msg.(type) {
	case string: // raw transport text, greetings, prompts, goodbyes ...
		return m

	case event.Failed:
		return failureText(m.Verb, string(m.Code))

	// ---- session -----------------------------------------------------------
	// The login events are consumed by conn.send before they ever reach the
	// queue; these cases exist so a stray one doesn't print a struct dump.

	case event.LoggedIn, event.PlayerCreated, event.LoginFailed, event.CreateFailed:
		return ""

	case event.LoggedOut:
		return m.Actor + " has logged out.\n"

	case event.Pong:
		return "Pong " + m.Target + ".\n"

	// ---- rooms -------------------------------------------------------------

	case event.RoomDescription:
		return renderRoom(m)

	case event.Entered:
		return m.Who + " enters.\n"

	case event.Left:
		// recall and other magical moves leave in no direction at all
		if m.Direction == direction.None {
			return m.Who + " leaves.\n"
		}
		return m.Who + " leaves " + strings.ToLower(m.Direction.String()) + ".\n"

	case event.Exits:
		return renderExits(m.Exits)

	// ---- objects -----------------------------------------------------------

	case event.Dropped:
		if m.Actor == self {
			return "Dropped.\n"
		}
		return m.Actor + " drops " + m.Item + ".\n"

	case event.Got:
		if m.Actor == self {
			return "Taken.\n"
		}
		return m.Actor + " gets " + m.Item + ".\n"

	case event.Equipped:
		return "Equipped.\n"

	case event.Worn:
		return "Done.\n"

	case event.Inventory:
		return renderInventory(m.Items)

	case event.Equipment:
		return renderEquipment(m.Items)

	// ---- talking -----------------------------------------------------------

	case event.Said:
		if m.Speaker == self {
			return "You say, \"" + m.Value + "\".\n"
		}
		return m.Speaker + " says, \"" + m.Value + "\".\n"

	case event.Told:
		if m.From == self {
			return "Ok.\n"
		}
		return m.From + " tells you, \"" + m.Value + "\".\n"

	case event.Shouted:
		if m.Speaker == self {
			return "Ok.\n"
		}
		return m.Speaker + " shouts, \"" + m.Value + "\".\n"

	// ---- the player --------------------------------------------------------

	case event.Who:
		return renderWho(m.Players)

	case event.Stat:
		return renderPlayerStat(m)

	// ---- combat ------------------------------------------------------------

	case event.Attacking:
		return "Ok.\n"

	case event.Struck:
		return renderViolence(self, m)

	case event.Died:
		return m.Target + " is dead!\n"

	case event.Restored:
		return m.Target + " is restored!\n"

	// ---- builder commands --------------------------------------------------

	case event.Loaded:
		return "Loaded.\n"

	case event.RoomStatus:
		return renderRoomStatus(m)

	default:
		log.Warn().Msgf("telnet render: no case for %T", msg)
		return fmt.Sprintf("%v\n", m)
	}
}

func renderEquipment(equipment []event.EquippedItem) string {
	if len(equipment) == 0 {
		return "Nothing equipped.\n"
	}
	var b strings.Builder
	b.WriteString("You are using:\n")
	for _, eq := range equipment {
		// include instance id just for testing, for now...
		b.WriteString(eq.Slot.String() + "\t" + eq.ShortDescription + "\t(" + eq.Id + ")\n")
	}
	b.WriteString("\n")
	return b.String()
}

func renderExits(exits []event.Exit) string {
	var b strings.Builder
	b.WriteString("Exits:\n")
	if len(exits) == 0 {
		b.WriteString("None!\n")
	} else {
		var exitStrs []string
		for _, exit := range exits {
			exitStrs = append(exitStrs, strings.ToLower(exit.Direction.String()))
		}
		b.WriteString(strings.Join(exitStrs, ", ") + "\n")
	}
	return b.String()
}

// renderInventory formats a list of inventory items as a string for display
// to a mud client.
func renderInventory(items []event.InventoryItem) string {
	if len(items) == 0 {
		return "You aren't carrying anything.\n"
	}
	var b strings.Builder
	b.WriteString("You are carrying:\n")
	for _, item := range items {
		b.WriteString("\t" + item.ShortDescription + "\n")
	}
	return b.String()
}

// renderPlayerStat formats a player's stats as a string for display to a mud
// client. The numbers are all zero until Phase 6 fills them in.
func renderPlayerStat(s event.Stat) string {
	abilities := rules.Abilities{
		Str: s.Strength,
		Dex: s.Dexterity,
		Con: s.Constitution,
		Int: s.Intelligence,
		Wis: s.Wisdom,
		Cha: s.Charisma,
	}
	var b strings.Builder
	b.WriteString("Status:\n")
	b.WriteString("Player:\t" + s.PlayerName + "\n")
	b.WriteString("Lineage:\t" + s.Lineage + "\tClass: " + s.Class + "\n")
	b.WriteString(fmt.Sprintf("Health:\t%d of %d\n", s.CurrentHealth, s.MaxHealth))
	b.WriteString("Location:\t" + player.NewLocation(s.ZoneId, s.RoomId).String() + "\n")
	b.WriteString("Abilities:\n")
	b.WriteString(fmt.Sprintf("\tStr: %d\t Dex: %d\t Con: %d\n", abilities.Str, abilities.Dex, abilities.Con))
	b.WriteString(fmt.Sprintf("\tWis: %d\t Int: %d\t Cha: %d\n", abilities.Wis, abilities.Int, abilities.Cha))
	b.WriteString("\n")
	return b.String()
}

// renderRoom formats a RoomDescription as the classic MUD room block:
// name, description, exits, then contents. Objects and mobs arrive
// as complete sentences (DescriptionOnGround / DescriptionInRoom) and
// print as-is; player names don't, so they get a verb here.
func renderRoom(rd event.RoomDescription) string {
	var b strings.Builder
	b.WriteString(rd.Name + "\n")
	if rd.Description != "" {
		b.WriteString(" " + rd.Description + "\n")
	}

	b.WriteString("[ Exits: " + rd.Exits + " ]\n")

	for _, o := range rd.Objects {
		b.WriteString(o + "\n")
	}
	for _, m := range rd.Mobs {
		b.WriteString(m + "\n")
	}
	for _, p := range rd.Players {
		b.WriteString(p + " is here.\n")
	}
	return b.String()
}

// renderRoomStatus is the builder's dump of everything in a room.
// Deliberately technical: the audience is someone editing content/.
func renderRoomStatus(rs event.RoomStatus) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Room %s.%s %q\n", rs.ZoneId, rs.Id, rs.Name)
	fmt.Fprintf(&b, "Zone: %s (%s)\n", rs.ZoneName, rs.ZoneId)
	if len(rs.Flags) > 0 {
		fmt.Fprintf(&b, "Flags: %s\n", strings.Join(rs.Flags, ", "))
	}
	for _, ex := range rs.Exits {
		fmt.Fprintf(&b, "  exit %-5s -> %s.%s\n", strings.ToLower(ex.Direction.String()), ex.ZoneId, ex.RoomId)
	}
	for _, p := range rs.Players {
		fmt.Fprintf(&b, "  player %s (%d/%d)\n", p.Name, p.CurrentHealth, p.MaxHealth)
	}
	for _, m := range rs.Mobs {
		fmt.Fprintf(&b, "  mob %s.%s %q (%d/%d) %s\n", m.ZoneId, m.DefinitionId, m.Name, m.CurrentHealth, m.MaxHealth, m.Id)
	}
	for _, i := range rs.Items {
		fmt.Fprintf(&b, "  obj %s.%s %q %s\n", i.ZoneId, i.DefinitionId, i.Name, i.Id)
	}
	return b.String()
}

func renderViolence(self string, s event.Struck) string {
	switch {
	case s.Attacker == self:
		if !s.Hit {
			return fmt.Sprintf("You miss %s.\n", s.Target)
		}
		return fmt.Sprintf("You hit %s for %d damage.\n", s.Target, s.Damage)
	case s.Target == self:
		if !s.Hit {
			return fmt.Sprintf("%s misses you.\n", s.Attacker)
		}
		return fmt.Sprintf("%s hits you for %d damage.\n", s.Attacker, s.Damage)
	default:
		if !s.Hit {
			return fmt.Sprintf("%s misses %s.\n", s.Attacker, s.Target)
		}
		// don't include damage numbers for the bystanders
		return fmt.Sprintf("%s hits %s.\n", s.Attacker, s.Target)
	}
}

func renderWho(players []event.WhoEntry) string {
	var b strings.Builder
	if len(players) == 0 {
		b.WriteString("There's nobody here!\n")
		return b.String()
	}
	b.WriteString("-- Who Is Here --\n")
	for _, p := range players {
		b.WriteString(p.PlayerName + " - " + p.RoomName + " - " + p.ZoneName + "\n")
	}
	return b.String()
}
