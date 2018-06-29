package combat

import "github.com/pkg/errors"

type FightLedger struct {
	fightMap map[Combatant]Combatant
}

func NewFightLedger() *FightLedger {
	return &FightLedger{
		fightMap: make(map[Combatant]Combatant),
	}
}

func (f *FightLedger) Fight(fighter Combatant, fightee Combatant) error {
	if f.IsFighting(fighter) {
		return errors.New("Fighter is already fighting someone")
	}
	f.fightMap[fighter] = fightee
	return nil
}

func (f *FightLedger) IsFighting(c Combatant) bool {
	_, exists := f.fightMap[c]
	return exists
}
