# Levels: power comes from gear

The progression loop: kill things near your power, they drop gear a little above it,
the gear raises your power, and bigger things become worth fighting. Crafting joins
later as another place gear comes from.

The model is Destiny's Power level (Diablo IV's Item Power and WoW's item level are
cousins of it): **your level is what you are wearing.** It follows from rules
CLAUDE.md already sets -- the gear is the truth, don't route mechanics through a
hidden number, derive combat stats from equipment.

## Decided

- **No XP, and no stored level.** Nothing earns points, and nothing on `player.Player`
  or `player.Record` holds a level. A cosmetic rank for `who` could come someday;
  nothing mechanical may read it.
- **Power is per item instance**, not per definition. The same rusty sword is power 5
  off a goblin and power 15 off an ogre, so one definition serves every tier. Like
  durability, it is set wherever the instance is made (the mob that dropped it, the
  zone that reset it, starting gear, later crafting) and saved as a `*int` in
  `InventoryRecord`, so records from before power existed still load.
- **A player's power is derived on every read, like role.** It's the average power of
  the equipped, unbroken items; **empty slots don't count.** Nothing equipped is
  power 0. No `Power` field on `Player`, for the same reason there is no `Role`
  field: a stored copy can disagree with the slots.
- **AC stays as it is**: armor type by slot, from `armor.json`. Plate means *harder to
  hit*; power means *stronger*. Two things, tuned separately.
- **Mobs get `"power"` in mobs.json**, defaulted by the loader the same way `ac` is.
  Zones get a power band in the manifest, which is how content says "this is the
  10-15 area".
- **Combat reads the power *difference*, clamped.** d20 vs AC has to stay bounded --
  if AC grew with power, by power 30 nothing would ever connect. So hit and damage
  stay as they are, and `attacker.Power - defender.Power` becomes a modifier on top.
  It's arithmetic, not a cutoff like role, so it doesn't break the rule against
  branching on a role.
- **Loot comes from loot tables on mobs**, goes into the corpse, and drops at the
  mob's power with a small chance of a bit more.
- **`consider <mob>`** compares powers. Considering a *player* is tabled.

## Tuning placeholders

Every number here is a guess until playtesting. They start as named constants in
`rules/`, each with a comment pointing back here, and move into a `rules/*.json`
together once there's something to tune against. Keep them in one place so that
move is one change.

| What | Placeholder | Where it matters |
|---|---|---|
| Regen interval | 5s (`mudtime.json`) | how long a grind stalls between fights |
| Regen amount | 5% of max health, at least 1 | same |
| Power delta clamp | ±10 | how far out of your league anything can be |
| To-hit per point of delta | +½ | whether a higher mob can be hit at all |
| Damage per point of delta | ±5% | how long fights above your power take |
| Loot power bump | 10% chance of +1, 2% of +2 | **the pace of the whole game** |

## Phases

Each one leaves `make check` green and is playable on its own.

### 1. Prerequisites

- **Regeneration.** Health comes back on a `regen` pulse for anyone not fighting,
  players and mobs both. Without it, dying (which leaves you on 1 hp) ends your
  progress for good, and so does a bad fight.
- **Weapon damage from the weapon.** `Player.WeaponDamageRoll` is a hardcoded `1d6`.
  Give `object.Definition` a damage roll from objects.json and read it from whatever
  is in `SlotWield`, bare hands otherwise.

*Done when:* a hurt player heals back while idle, and swapping weapons changes the
damage.

### 2. Power

- `object.Instance.Power`, saved with the record the way durability is.
- `Equipment.Power()`: the average over equipped, unbroken items.
- `mobile.Definition.Power` from mobs.json, and a power band per zone in the manifest.
- The combat modifier from the clamped delta, in `combat/` against a new
  `Combatant.Power()`.
- Power shown in `stat` and `equipment`; `consider <mob>`.

*Done when:* the same fight goes differently as your gear gets better, and `consider`
says so before you start it.

### 3. Loot

- Loot tables in mobs.json: `"loot": [{"object": "wrathrock/rusty_sword", "chance": 25}]`.
- Corpses hold the drops, which finishes the `TODO transfer m's possessions` in
  `world/corpse.go`. Drop power is the mob's, with the bump from the table above.
- `get <item> from <container>`, and `get all from corpse`.
- Corpses decay, or rooms fill up with them.

*Done when:* killing something gives you something worth wearing.

### 4. Content

Power bands on the zones and mobs to fill them. Content work, not code, and where the
placeholders above get replaced.

### 5. Crafting (later)

Another place an item's power comes from, so it slots into the same model:

- **salvage**: break gear down into materials, at that gear's tier
- **craft**: make gear, at the materials' tier
- **infuse**: feed a stronger item into a favourite to raise its power (Destiny's trick)
- **repair**: the durability sink. `object.Instance.Repair` exists and nothing calls it.

## Tabled

- **Considering a player** -- their gear, role, power. PvP isn't on the board.
- **Abilities from gear**: a ring of healing gives a healing cast, a mace of smackdown
  gives a big hit with a stun. It's the natural next step after power: gear decides
  what you can *do*, as well as how strong you are.

## Known risk in "average only what you wear"

Wearing one power-20 ring and nothing else makes you power 20, the same as a full
power-20 kit. That's fine while power only shifts combat numbers a little, and it may
matter less once abilities come from gear. If it becomes a problem, two options that
keep the "empty slots don't count" feel:

- **weight by slot**: the weapon and body count for more than a ring
- **split by side**: offense from the weapon's power, defense from the armor's average
