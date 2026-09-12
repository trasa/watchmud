package spaces

import (
	"github.com/trasa/watchmud/direction"
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
