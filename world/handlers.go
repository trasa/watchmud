package world

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud/command"
	"github.com/trasa/watchmud/event"
	"github.com/trasa/watchmud/gameserver"
)

// HandleIncomingMessage runs the handler for one command, then persists the
// player.
//
// The switch is the whole dispatch table: no string keys, no reflection, and
// each handler is handed its command already typed.
func (w *World) HandleIncomingMessage(msg *gameserver.HandlerParameter) error {
	switch cmd := msg.Command.(type) {
	case command.Drop:
		w.handleDrop(msg, cmd)
	case command.Equip:
		w.handleEquip(msg, cmd)
	case command.Exits:
		w.handleExits(msg, cmd)
	case command.Get:
		w.handleGet(msg, cmd)
	case command.Inventory:
		w.handleInventory(msg, cmd)
	case command.Kill:
		w.handleKill(msg, cmd)
	case command.Load:
		w.handleLoad(msg, cmd)
	case command.Logout:
		w.handleLogout(msg, cmd)
	case command.Look:
		w.handleLook(msg, cmd)
	case command.Move:
		w.handleMove(msg, cmd)
	case command.Ping:
		w.handlePing(msg, cmd)
	case command.Recall:
		w.handleRecall(msg, cmd)
	case command.Remove:
		w.handleRemove(msg, cmd)
	case command.Restore:
		w.handleRestore(msg, cmd)
	case command.Role:
		w.handleRole(msg, cmd)
	case command.RoomStatus:
		w.handleRoomStatus(msg, cmd)
	case command.Say:
		w.handleSay(msg, cmd)
	case command.ShowEquipment:
		w.handleShowEquipment(msg, cmd)
	case command.Stat:
		w.handleStat(msg, cmd)
	case command.Tell:
		w.handleTell(msg, cmd)
	case command.TellAll:
		w.handleTellAll(msg, cmd)
	case command.Wear:
		w.handleWear(msg, cmd)
	case command.Who:
		w.handleWho(msg, cmd)
	default:
		log.Warn().Msgf("world.HandleIncomingMessage: UNHANDLED command %T", msg.Command)
		msg.Fail(event.UnknownCommand)
		return fmt.Errorf("unhandled command %T", msg.Command)
	}
	return w.save(msg)
}

// save persists the player after a handler ran.
// TODO what if the player has changed some other player somehow?
// (stabbed them, stole from them, etc.) See #32.
// Saving after every single message is stupid but free against a map; the
// Store interface means a real implementation can batch or debounce without
// world/ knowing.
func (w *World) save(msg *gameserver.HandlerParameter) error {
	if msg.Player == nil {
		return nil
	}
	return w.store.Save(msg.Player.Record())
}
