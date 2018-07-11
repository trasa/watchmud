package world

import (
	"errors"
	"fmt"
	"github.com/trasa/watchmud-message/direction"
	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/spaces"
	"log"
	"time"
)

// Walk through all the mob instances that are in this world
// right now and tell them all to do something, if they have
// anything they want to do.
func (w *World) DoMobileActivity() {
	// for each mob in the world
	// wake it up and tell it to do stuff
	// don't limit this to per zone or per room or something
	// remember that mobs can leave the zone they started out in
	// (if programmed to)
	// (or if they really really want to)
	for _, mob := range w.mobileRooms.GetAllMobiles() {
		if mob.CanWander() {
			switch mob.Definition.Wandering.Style {
			case mobile.WANDER_RANDOM:
				// do random wander within the zone
				if err := w.doMobRandomWander(mob); err != nil {
					log.Printf("World.DoMobileActivity: %s error randomly wandering: %s", mob.Definition.Id, err)
				}
			case mobile.WANDER_FOLLOW_PATH:
				if err := w.doMobFollowPathWander(mob); err != nil {
					log.Printf("World.DoMobileActivity: %s error following path: %s", mob.Definition.Id, err)
				}
			default:
				// unknown or unhandled wandering style, do nothing.
			}
		}
	}
}

// pick a direction that is within the mob's zone and walk to it, if possible.
func (w *World) doMobRandomWander(mob *mobile.Instance) error {
	mob.LastWanderingTime = time.Now()
	mobRoom := w.getRoomContainingMobile(mob)
	if mobRoom == nil {
		return errors.New(fmt.Sprintf("Mobile ID '%s' can't randomly wander - not in a room at all!", mob.Definition.Id))
	}
	// test wandering percentage
	if mob.CheckWanderChance() {
		dir := mobRoom.PickRandomDirection(true)
		if dir == direction.NONE {
			return errors.New(fmt.Sprintf("Mobile ID '%s' is in a room without exit and can't wander out of it.", mob.Definition.Id))
		}
		w.moveMobile(mob, dir, mobRoom, mobRoom.Get(dir))
		//log.Printf("World.doMobRandomWander: %s randomly wanders to %s", mob.Definition.Id, mobRoom.Get(dir))
	}
	return nil
}

func (w *World) doMobFollowPathWander(mob *mobile.Instance) error {
	mob.LastWanderingTime = time.Now()
	mobRoom := w.getRoomContainingMobile(mob)
	if mobRoom == nil {
		return errors.New(fmt.Sprintf("Mobile ID '%s' can't follow path - not in a room at all!", mob.Definition.Id))
	}
	if mob.CheckWanderChance() {
		dir, changeDirection, err := getNextDirectionOnPath(mob, mobRoom)
		if err != nil {
			return err
		}
		if dir == direction.NONE {
			return errors.New(fmt.Sprintf("doMobFollowPathWander: mobile ID '%s' can't figure out next place to go to (current '%s', path '%s')",
				mob.Definition.Id, mobRoom.Id, mob.Definition.Wandering.Path))
		}
		if changeDirection {
			mob.WanderingForward = !mob.WanderingForward
		}
		w.moveMobile(mob, dir, mobRoom, mobRoom.Get(dir))
		//log.Printf("World.doMobFollowPathWander: %s moves to %s", mob.Definition.Id, mobRoom.Get(dir))
	}
	return nil
}

// Determine what direction this mob should travel next to stay on its path.
// Takes mob.WanderingForward into account and will reverse index at the path bounaries,
// returning changeDirection=true in that case, but WILL NOT update the state or
// modify the mob or room instances in any way.
func getNextDirectionOnPath(mob *mobile.Instance, mobRoom *spaces.Room) (dir direction.Direction, changeDirection bool, err error) {
	currentIndex, err := mob.GetIndexOnPath(mobRoom.Id)
	if err != nil {
		return direction.NONE, false, err
	}
	nextIndex := -1

	if currentIndex < 0 {
		// note: this might be OK (if the mob was pulled off the path for some reason?)
		// TODO should it change to a random walk? or just wait here, or?
		return direction.NONE, false, errors.New(fmt.Sprintf("Couldn't find current room %s in wander path %s", mobRoom.Id, mob.Definition.Wandering.Path))
	}
	if mob.WanderingForward {
		nextIndex = currentIndex + 1
		// is it time to change directions?
		if nextIndex > len(mob.Definition.Wandering.Path)-1 {
			// yep
			changeDirection = true
			nextIndex = currentIndex - 1
		}
	} else {
		nextIndex = currentIndex - 1
		// is it time to change directions?
		if nextIndex < 0 {
			// yep
			changeDirection = true
			nextIndex = currentIndex + 1
		}
	}
	roomToFind := mob.Definition.Wandering.Path[nextIndex]
	for _, rexit := range mobRoom.GetRoomExits(false) {
		if rexit.Room.Id == roomToFind {
			dir = rexit.Direction
			break
		}
	}
	if dir == direction.NONE {
		return direction.NONE, false, errors.New(fmt.Sprintf("Couldn't find destination room %s from current room exits %v", roomToFind, mobRoom.GetRoomExits(false)))
	}
	return
}
