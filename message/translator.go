package message

import (
	"errors"
	"github.com/trasa/watchmud/direction"
	"github.com/trasa/watchmud/slot"
	"log"
	"strings"
)

// Parse line into tokens
func Tokenize(line string) []string {
	return strings.Fields(line)
}

// Turn a request of format type "line" into a Game Message
//
// have an input string like
// tell bob hi there
// turn into Message.Inner = message.TellRequest { "bob", "hi there" }
// and so on.
// note that not all commands can be parsed from line input.
func TranslateLineToMessage(tokens []string) (msg *GameMessage, err error) {

	var payload interface{}
	// first parse into tokens
	if len(tokens) == 0 {
		log.Printf("no tokens")
		payload = Ping{Target: "x"}
		log.Printf("%v", payload)
	} else {
		switch tokens[0] {
		case "":
			// No Op
			payload = Ping{}

		case "drop":
			payload = DropRequest{
				Target: tokens[1],
			}

		case "exits", "exit", "ex":
			payload = ExitsRequest{}

		case "get":
			payload = GetRequest{
				Targets: tokens[1:],
			}

		case "inv", "inventory":
			payload = InventoryRequest{}

		case "look", "l":
			payload = LookRequest{
				ValueList: tokens[1:],
			}

		case "n", "north", "s", "south", "e", "east", "w", "west", "u", "up", "d", "down":
			var d direction.Direction
			if d, err = direction.StringToDirection(tokens[0]); err == nil {
				payload = MoveRequest{
					Direction: int32(d),
				}
			}

		case "'", "say":
			if len(tokens) >= 2 {
				payload = SayRequest{
					Value: strings.Join(tokens[1:], " "),
				}
			} else {
				err = errors.New("What do you want to say?")
			}

		case "tell", "t":
			if len(tokens) >= 3 {
				payload = TellRequest{
					ReceiverPlayerName: tokens[1],
					Value:              strings.Join(tokens[2:], " "),
				}
			} else {
				// some sort of error about malformed tell request...
				err = errors.New("usage: tell [somebody] [something]")
			}

		case "tellall", "ta":
			if len(tokens) >= 2 {
				payload = TellAllRequest{
					Value: strings.Join(tokens[1:], " "),
				}
			} else {
				err = errors.New("usage: tellall [something]")
			}

		case "wield":
			if len(tokens) >= 2 {
				payload = EquipRequest{
					Target:       tokens[1],
					SlotLocation: int32(slot.Wield),
				}
			} else {
				err = errors.New("What do you want to wield?")
			}

		case "who":
			payload = WhoRequest{}

		default:
			err = errors.New("Unknown request: " + tokens[0])
		}
	}
	if err == nil {
		msg, err = NewGameMessage(payload)
	}
	return
}
