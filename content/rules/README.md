# These Are The Rules

## Armor

Describes the AC bonuses for each type of armor depending on which slot it's in.

The same number does double duty: it is the AC the piece adds, and it is what the
piece contributes to any role marked `from_armor` in roles.json. Armor that protects
you more is armor that argues harder that you're the one standing in front, so there
is nothing to keep in sync.

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