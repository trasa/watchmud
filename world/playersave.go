package world

import (
	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud/player"
)

func (w *World) QueuePlayerRecords() {
	for p := range w.Players() {
		if err := w.store.Save(w.record(p)); err != nil {
			log.Error().Err(err).Msgf("queuePlayerRecords : %v", err)
		}
	}
}

// Record is the player's record plus where they are standing. The player
// can't fill that in themselves: location lives in playerToRoom, not on the
// player. Every save goes through here, or the room is forgotten.
func (w *World) record(p *player.Player) *player.Record {
	rec := p.Record()
	if r := w.playerToRoom.Get(p); r != nil {
		rec.LastZoneId = r.Zone.Id
		rec.LastRoomId = r.Id
	}
	return rec
}
