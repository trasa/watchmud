package playergenerator

import (
	"github.com/trasa/watchmud/db"
)

type RawPlayer struct {}

func generate(race db.RaceData, class db.ClassData) RawPlayer {

	// TODO take racial attrs and class preferences and put together
	// the set of abilities to start with, then send that to the
	// player so they can make changes.

	/* ability scores determined by class:
	15, 14, 13, 12, 10, 8
	 */
	return RawPlayer{}
}
