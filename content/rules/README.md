# These Are The Rules

## Armor

Describes the AC bonuses for each type of armor depending on which slot it's in.

The same number does double duty: it is the AC the piece adds, and it is what the
piece contributes to any role marked `from_armor` in roles.json. Armor that protects
you more is armor that argues harder that you're the one standing in front, so there
is nothing to keep in sync.

AC starts at 10 (`rules.BaseArmorClass`) and a hit is a d20 that *meets or beats* it,
so 10 is unarmored and hit 55% of the time. A mob's `"ac"` in mobs.json is on that same
absolute scale -- 10 is an unarmored creature, not a bonus on top of one -- and a mob
that doesn't name one gets 10. Writing `"ac": 0` means a creature that cannot be missed,
since a d20 never rolls below 1.

## Roles

This is the replacement for classes: your role is determined by what gear you're wearing.
Healer Gear = you're a healer, etc.

A role with `"from_armor": true` is one that armor argues for on its own, using the
table above -- that's Tank. Everything else says what it's worth by hand, with a
`"roles"` key on the object: a knife has no armor type to derive anything from.

## Species

This is the definition of species and lineage, which replaces 'Race' and 'Class' and
is now mostly cosmetic. Your choice of lineage doesn't dictate any stats, just how
you feel like role-playing.

## Starting Gear

What a brand-new character is created holding: `starting_gear.json`, a list of object
definitions named by the zone that defines them, in the order they should arrive.

```json
{ "zone": "wrathrock", "object": "training_dagger", "equip": true }
```

`"equip": true` means the character starts with it worn; leave it off and the item is
only carried. There is no slot here on purpose -- the object definition already names
the one slot it goes in, and a second copy of that would only be something to keep in
step. An item marked `equip` that isn't wearable, two items claiming the same slot, or
an object or zone that doesn't exist, all fail startup rather than quietly handing out
less than the file says.

The file is optional: no `starting_gear.json` means new characters start with nothing.

Nothing in here says what role the kit adds up to, because nothing can -- a role is read
off the equipment at the moment it's asked for. The wrathrock kit (a training dagger, a
wool tunic and a cloak, plus a waterskin) makes a level 1 character a Striker with AC 11,
and stops mattering the moment they wear something else.
