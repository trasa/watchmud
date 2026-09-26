package world

import "github.com/watchmud/watchmud/rules"

// Regenerate gives everyone who isn't fighting a little health back:
// players and mobs both, or a boss could be worn down by dying at it
// over and over.
//
// Not while fighting, on either side of the fight. Healing mid-fight
// turns a slow fight into one that never ends. And never the dead:
// regen is not resurrection.
func (w *World) Regenerate() {
	for p := range w.Players() {
		if p.Dead() || w.fightLedger.InFight(p) {
			continue
		}
		p.RestoreHealth(rules.RegenAmount(p.MaxHealth()))
	}
	for _, mob := range w.Mobiles() {
		if mob.Dead() || w.fightLedger.InFight(mob) {
			continue
		}
		mob.RestoreHealth(rules.RegenAmount(mob.Definition.MaxHealth))
	}
}
