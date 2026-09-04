package spaces

import (
	"sort"

	"github.com/trasa/watchmud-message/direction"
)

// RoomExit is a direction to another room.
type RoomExit struct {
	Direction direction.Direction
	Room      *Room
}

type roomExitHolder struct {
	dirs []RoomExit
}

func (re roomExitHolder) Len() int           { return len(re.dirs) }
func (re roomExitHolder) Less(i, j int) bool { return re.dirs[i].Direction < re.dirs[j].Direction }
func (re roomExitHolder) Swap(i, j int)      { re.dirs[i], re.dirs[j] = re.dirs[j], re.dirs[i] }

// GetRoomExits returns the exits from this room.
// Uses the direction.Direction ordering.
// Does not take locks, doors, closures, etc. into account.
func (r *Room) GetRoomExits(limitToZone bool) []RoomExit {
	holder := roomExitHolder{}
	for dir, dest := range r.directions {
		if !limitToZone || r.Zone == dest.Zone {
			holder.dirs = append(holder.dirs, RoomExit{dir, dest})
		}
	}
	sort.Sort(holder)
	return holder.dirs
}
