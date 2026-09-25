package mobile

import "github.com/watchmud/watchmud/object"

// LootEntry is one line of a mob's loot table: what might drop, and the
// percent chance it does. Each line is rolled on its own when the mob dies,
// so a mob can drop everything on its table or nothing. See LEVELS.md.
type LootEntry struct {
	Object *object.Definition
	Chance int
}
