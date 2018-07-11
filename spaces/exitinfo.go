package spaces

import (
	"github.com/trasa/watchmud-message/direction"
	"log"
	"math/rand"
	"strings"
	"time"
)

// Get all the valid exits from this room.
// Note the ordering of the letters in the exit string is important!
// For example, with exits north, south, and up,
// this returns "nsu"
// TODO: exits can be locked and/or closed, this doesn't handle that.
func (r *Room) GetExitString() string {
	exits := []string{}

	for _, rexit := range r.GetRoomExits(false) {
		if s, err := direction.DirectionToAbbreviation(rexit.Direction); err == nil {
			exits = append(exits, s)
		}
	}
	return strings.Join(exits, "")
}

// Is there a valid exit in this direction in this room?
// TODO what about exits that are locked or closed?
func (r *Room) HasExit(dir direction.Direction) bool {
	_, ok := r.directions[dir]
	return ok
}

// Get the exit for this direction. Will return nil if there
// isn't a valid exit that way.
// TODO what about exits that are locked or closed?
func (r *Room) Get(dir direction.Direction) (dest *Room) {
	return r.directions[dir]
}

func (r *Room) Set(dir direction.Direction, destRoom *Room) {
	r.directions[dir] = destRoom
}

// Of the directions available for travel (could be locked, closed...)
// pick one of them. If there aren't any, return none.
func (r *Room) PickRandomDirection(limitToZone bool) direction.Direction {
	exits := r.GetRoomExits(limitToZone)
	if len(exits) == 0 {
		return direction.NONE
	} else {
		desired := rand.New(rand.NewSource(time.Now().Unix())).Int31n(int32(len(exits)))
		// iterate to the ith member of exits
		i := int32(0)
		for _, re := range exits {
			if i == desired {
				return re.Direction
			}
			i++
		}
		// inconceivable!
		log.Printf("Room.PickRandomDirection: Bizzare RandomDirection picked. len=%d, desired=%d", len(exits), desired)
		return direction.NONE
	}
}
