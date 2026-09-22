package rules

// DurabilityTable says what a piece of gear starts out able to take, from
// content/rules/durability.json. Plate outlasts cloth; everything that isn't
// armor falls back on one number.
//
// It is a table rather than a field on every object for the same reason
// armor.json is: a builder tunes what leather is worth once, and nobody keeps
// a second number in step with it. An individual object can still override it
// -- a dwarven blade is allowed to be special -- but it has to say so.
type DurabilityTable struct {
	// Default is what gear with no armor type starts at: a weapon, a censer,
	// a ring. Zero means nothing wears out, which is what an absent
	// durability.json gets you.
	Default int `json:"default"`

	// Armor is by armor type, since that is what says how much punishment a
	// piece is built to take. A type the table doesn't mention falls back to
	// Default.
	Armor map[ArmorType]int `json:"armor"`

	// OnDeathPercent is what dying costs every piece you died in, as a
	// percentage of what that piece started at.
	//
	// A percentage rather than a flat number of points so that dying in plate
	// and dying in a wool tunic cost the same number of *deaths* rather than
	// the same number of points -- eight off an eighty-point breastplate and
	// two off a twenty-point tunic are the same fraction of a life's worth of
	// gear. Zero is death costing nothing, which is what content that says
	// nothing gets.
	OnDeathPercent int `json:"on_death_percent"`
}

// DeathLoss is how many points dying takes off one piece that started at
// maxDurability.
//
// Always at least one point when the rule is switched on at all: a percentage
// that rounds to nothing would make cheap gear immortal in exactly the
// situation the rule exists to punish.
func (t DurabilityTable) DeathLoss(maxDurability int) int {
	if t.OnDeathPercent <= 0 || maxDurability <= Indestructible {
		return 0
	}
	return max(1, maxDurability*t.OnDeathPercent/100)
}

// MaxFor is what a new piece of gear of this armor type starts at.
func (t DurabilityTable) MaxFor(armorType ArmorType) int {
	if armorType != ArmorTypeNone {
		if max, listed := t.Armor[armorType]; listed {
			return max
		}
	}
	return t.Default
}

// Indestructible is what a max durability of zero means: gear that never
// wears out. It reads as the absence of the mechanic rather than as gear that
// is born broken, which is the only sane way round -- content that says
// nothing about durability should hand out working equipment, and a server
// with no durability.json at all should play exactly as it did before there
// was one.
const Indestructible = 0
