package combat

import "github.com/pkg/errors"

type FightLedger struct {
	fightMap map[Combatant]*Fight
}

func NewFightLedger() *FightLedger {
	return &FightLedger{
		fightMap: make(map[Combatant]*Fight),
	}
}

func (f *FightLedger) Fight(fighter Combatant, fightee Combatant) error {
	if f.IsFighting(fighter) {
		return errors.New("Fighter is already fighting someone")
	}
	f.fightMap[fighter] = NewFight(fighter, fightee)
	return nil
}

func (f *FightLedger) IsFighting(c Combatant) bool {
	_, exists := f.fightMap[c]
	return exists
}

func (f *FightLedger) GetFight(fighter Combatant) *Fight {
	return f.fightMap[fighter]
}

func (f *FightLedger) GetFights() (result []*Fight) {
	for _, v := range f.fightMap {
		result = append(result, v)
	}
	return result
}