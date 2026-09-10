package telnet

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	message "github.com/trasa/watchmud-message"
)

// Render turns anything sent to a connection into the text a telnet client
// sees. It is the single checkpoint between the game's message vocabulary
// and the wire, so keep game logic out of it.
func render(msg any) string {
	// TODO (phase 5) we should be switching on *message.SomeResponse and not message.SomeResponse
	// but that doesn't cause problems currently (protobufs isn't being used as a transport)
	// and we'll have to deal with it later.

	switch m := msg.(type) {
	case string: // raw transport text, greetings, prompts, goodbyes ...
		return m

	//case message.GetResponse:

	//case message.EnterRoomNotification:

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

	case message.RecallResponse:
		if !m.Success {
			return failureText(m.ResultCode)
		}
		return renderRoom(m.RoomDescription)

	case message.RoomDescription:
		// if !m.Success -- wait, this doesn't declare a Success method?! ugh.
		return renderRoom(&m)

	case message.SayNotification:
		if !m.Success {
			return failureText(m.ResultCode)
		}
		return m.Sender + " says, \"" + m.Value + "\"\n"

	case message.SayResponse:
		if !m.Success {
			return failureText(m.ResultCode)
		}
		return "you say, \"" + m.Value + "\"\n"

	case message.TellAllNotification:
		if !m.Success {
			return failureText(m.ResultCode)
		}
		return m.Sender + " shouts, \"" + m.Value + "\"\n"

	case message.TellNotification:
		if !m.Success {
			return failureText(m.ResultCode)
		}
		return m.Sender + " tells you, \"" + m.Value + "\"\n"

	default:
		log.Warn().Msgf("telnet render: no case for %T", msg)
		return fmt.Sprintf("%v\n", m)
	}
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

func failureText(resultCode string) string {
	return fmt.Sprintf("Failure: %v\n", resultCode)
}
