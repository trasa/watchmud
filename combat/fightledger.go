package combat

import (
	"fmt"
	"maps"
	"slices"
	"uuid"

	"github.com/rs/zerolog/log"
)

type FightLedger struct {
	fightMap map[uuid.UUID]*Fight
}

func NewFightLedger() *FightLedger {
	return &FightLedger{
		fightMap: make(map[uuid.UUID]*Fight),
	}
}

func (f *FightLedger) Fight(fighter Combatant, fightee Combatant, zoneId string, roomId string) error {
	if f.IsFighting(fighter) {
		// TODO fixme
		return fmt.Errorf("fighter is already fighting someone")
	}
	f.fightMap[fighter.Id()] = newFight(fighter, fightee, zoneId, roomId)

	if !f.IsFighting(fightee) {
		f.fightMap[fightee.Id()] = newFight(fightee, fighter, zoneId, roomId)
	}
	return nil
}

func (f *FightLedger) IsFighting(c Combatant) bool {
	allFights := slices.Collect(maps.Values(f.fightMap))
	log.Debug().Msgf("all fights: %v", allFights)
	_, exists := f.fightMap[c.Id()]
	return exists
}

func (f *FightLedger) IsBeingFought(c Combatant) bool {
	for _, fight := range f.fightMap {
		if fight.Fightee == c {
			return true
		}
	}
	return false
}

func (f *FightLedger) GetFight(fighter Combatant) *Fight {
	return f.fightMap[fighter.Id()]
}

func (f *FightLedger) GetFights() (result []*Fight) {
	for _, v := range f.fightMap {
		result = append(result, v)
	}
	return result
}

func (f *FightLedger) EndFight(fighter Combatant) {
	delete(f.fightMap, fighter.Id())
}
