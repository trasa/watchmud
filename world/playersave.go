package world

import (
	"github.com/rs/zerolog/log"
	"github.com/watchmud/watchmud/player"
)

// QueuePlayerRecords hands a record of everyone in the world to the store.
// This is the timed save, and the last one at shutdown: whatever changed a
// player -- a command, a fight, a death -- is in the next one.
func (w *World) QueuePlayerRecords() {
	for p := range w.Players() {
		if err := w.store.Save(w.record(p)); err != nil {
			log.Error().Err(err).Str("player", p.Name()).Msg("queueing player save")
		}
	}
}

// Record is the player's record plus where they are standing. The player
// can't fill that in themselves: location lives in the world's Occupancy, not on the
// player. Every save goes through here, or the room is forgotten.
func (w *World) record(p *player.Player) *player.Record {
	rec := p.Record()
	// don't use version that returns Void if they aren't in a room,
	// we want to test for that case.
	if r := w.occupancy.RoomOfPlayer(p); r != nil {
		rec.LastZoneId = r.Zone.Id
		rec.LastRoomId = r.Id
	}
	return rec
}
