package telnet

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	message "github.com/trasa/watchmud-message"
	"github.com/trasa/watchmud-message/direction"
	"github.com/trasa/watchmud-message/slot"
)

// Render turns anything sent to a connection into the text a telnet client
// sees. It is the single checkpoint between the game's message vocabulary
// and the wire, so keep game logic out of it.
func render(msg any, self string) string {
	// TODO (phase 5) we should be switching on *message.SomeResponse and not message.SomeResponse
	// but that doesn't cause problems currently (protobufs isn't being used as a transport)
	// and we'll have to deal with it later.

	switch m := msg.(type) {
	case string: // raw transport text, greetings, prompts, goodbyes ...
		return m

	case message.DeathNotification:
		return m.Target + " is dead!\n"

	case message.DropNotification:
		if !m.Success {
			return failureText(m.ResultCode)
		}
		return m.PlayerName + " drops " + m.Target + ".\n"

	case message.DropResponse:
		if !m.Success {
			return failureText(m.ResultCode)
		}
		return "Dropped.\n"

	case message.EnterRoomNotification:
		if !m.Success {
			return failureText(m.ResultCode)
		}
		return m.Name + " enters.\n"

	case message.EquipResponse:
		if !m.Success {
			return failureText(m.ResultCode)
		}
		return "Equipped.\n"

	case message.ExitsResponse:
		if !m.Success {
			return failureText(m.ResultCode)
		}
		return renderExits(m.ExitInfo)

	case message.GetNotification:
		if !m.Success {
			return failureText(m.ResultCode)
		}
		return m.PlayerName + " gets " + m.Target + ".\n"

	case message.GetResponse:
		if !m.Success {
			return failureText(m.ResultCode)
		}
		return "Taken.\n"

	case message.InventoryResponse:
		if !m.Success {
			return failureText(m.ResultCode)
		}
		return renderInventory(m.InventoryItems)

	case message.KillResponse:
		if !m.Success {
			return failureText(m.ResultCode)
		}
		return "Ok.\n"

	case message.LeaveRoomNotification:
		if !m.Success {
			return failureText(m.ResultCode)
		}
		return m.Name + " leaves " + strings.ToLower(direction.Direction(m.Direction).String()) + ".\n"

	case message.LoadResponse:
		if !m.Success {
			return failureText(m.ResultCode)
		}
		return "Loaded.\n"

	case message.LogoutResponse:
		if !m.Success {
			return failureText(m.ResultCode)
		}
		return "Bye.\n"

	case message.LookNotification:
		if !m.Success {
			return failureText(m.ResultCode)
		}
		return renderRoom(m.RoomDescription)

	case message.LookResponse:
		if !m.Success {
			return failureText(m.ResultCode)
		}
		return renderRoom(m.RoomDescription)

	case message.MoveResponse:
		if !m.Success {
			return failureText(m.ResultCode)
		}
		return renderRoom(m.RoomDescription)

	case message.ShowEquipmentResponse:
		if !m.Success {
			return failureText(m.ResultCode)
		}
		return renderEquipment(m.EquipmentInfo)

	case message.RecallResponse:
		if !m.Success {
			return failureText(m.ResultCode)
		}
		return renderRoom(m.RoomDescription)

	case message.RestoreNotification:
		return m.Target + " is restored!\n"

	case message.RestoreResponse:
		if !m.Success {
			return failureText(m.ResultCode)
		}
		return "Restored.\n"

	case message.RoomDescription:
		// if !m.Success -- wait, this doesn't declare a Success method?! ugh.
		return renderRoom(&m)

	case message.SayNotification:
		if !m.Success {
			return failureText(m.ResultCode)
		}
		return m.Sender + " says, \"" + m.Value + "\".\n"

	case message.SayResponse:
		if !m.Success {
			return failureText(m.ResultCode)
		}
		return "You say, \"" + m.Value + "\".\n"

	case message.TellAllNotification:
		if !m.Success {
			return failureText(m.ResultCode)
		}
		return m.Sender + " shouts, \"" + m.Value + "\".\n"

	case message.TellNotification:
		if !m.Success {
			return failureText(m.ResultCode)
		}
		return m.Sender + " tells you, \"" + m.Value + "\".\n"

	case message.TellAllResponse:
		if !m.Success {
			return failureText(m.ResultCode)
		}
		return "Told.\n"

	case message.TellResponse:
		if !m.Success {
			return failureText(m.ResultCode)
		}
		return "Told.\n"

	case message.ViolenceNotification:
		return renderViolence(self, m.SuccessfulHit, m.Fighter, m.Fightee, m.Damage)

	case message.WearResponse:
		if !m.Success {
			return failureText(m.ResultCode)
		}
		return "Done.\n"

	case message.WhoResponse:
		if !m.Success {
			return failureText(m.ResultCode)
		}
		return renderWho(m.PlayerInfo)

	default:
		log.Warn().Msgf("telnet render: no case for %T", msg)
		return fmt.Sprintf("%v\n", m)
	}
}

func renderEquipment(equipment []*message.ShowEquipmentResponse_EquipmentInfo) string {
	if len(equipment) == 0 {
		return "Nothing equipped.\n"
	}
	var b strings.Builder
	b.WriteString("You are using:\n")
	for _, eq := range equipment {
		l := slot.Location(eq.SlotLocation).String()
		// include instance id just for testing, for now...
		b.WriteString(l + "\t" + eq.ShortDescription + "\t(" + eq.Id + ")\n")
	}
	b.WriteString("\n")
	return b.String()
}

func renderExits(exits []*message.ExitInfo) string {
	var b strings.Builder
	b.WriteString("Exits:\n")
	if len(exits) == 0 {
		b.WriteString("None!\n")
	} else {
		var exitStrs []string
		for _, exit := range exits {
			exitStrs = append(exitStrs, strings.ToLower(direction.Direction(exit.Direction).String()))
		}
		b.WriteString(strings.Join(exitStrs, ", ") + "\n")
	}
	return b.String()
}

// renderInventory formats a list of inventory items as a string for display
// to a mud client.
func renderInventory(items []*message.InventoryResponse_InventoryItem) string {
	var b strings.Builder
	b.WriteString("You are carrying:\n")
	for _, item := range items {
		b.WriteString(item.ShortDescription + "\n")
	}
	return b.String()
}

// renderRoom formats a RoomDescription as the classic MUD room block:
// name, description, exits, then contents. Objects and mobs arrive
// as complete sentences (DescriptionOnGround / DescriptionInRoom) and
// print as-is; player names don't, so they get a verb here.
func renderRoom(rd *message.RoomDescription) string {
	if rd == nil {
		return "You can't see anything.\n"
	}
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

func renderViolence(self string, success bool, attacker string, target string, damage int32) string {
	switch {
	case attacker == self:
		if !success {
			return fmt.Sprintf("You miss %s.\n", target)
		}
		return fmt.Sprintf("You hit %s for %d damage.\n", target, damage)
	case target == self:
		if !success {
			return fmt.Sprintf("%s misses you.\n", attacker)
		}
		return fmt.Sprintf("%s hits you for %d damage.\n", attacker, damage)
	default:
		if !success {
			return fmt.Sprintf("%s misses %s.\n", attacker, target)
		}
		// don't include damage numbers for the bystanders
		return fmt.Sprintf("%s hits %s.\n", attacker, target)
	}
}

func renderWho(players []*message.WhoResponse_PlayerInfo) string {
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
