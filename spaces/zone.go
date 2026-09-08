package spaces

import (
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/zonereset"
)

type Zone struct {
	Id                string
	Rooms             map[string]*Room              // id -> *room
	ObjectDefinitions map[string]*object.Definition // id (no zone) -> object.Definition
	MobileDefinitions map[string]*mobile.Definition // id -> mobile.Definition
	Name              string
	Commands          []ZoneCommand
	ResetMode         zonereset.Mode
	LastReset         time.Time
	Lifetime          time.Duration
}

func NewZone(id string, name string, resetMode zonereset.Mode, lifetime time.Duration) *Zone {
	return &Zone{
		Id:                id,
		Name:              name,
		Rooms:             make(map[string]*Room),
		ObjectDefinitions: make(map[string]*object.Definition),
		MobileDefinitions: make(map[string]*mobile.Definition),
		ResetMode:         resetMode,
		Lifetime:          lifetime,
	}
}

func (z *Zone) AddRoom(r *Room) {
	r.Zone = z
	z.Rooms[r.Id] = r
}

func (z *Zone) AddObjectDefinition(obj *object.Definition) {
	obj.ZoneId = z.Id
	z.ObjectDefinitions[obj.Identifier()] = obj
}

func (z *Zone) AddMobileDefinition(mob *mobile.Definition) {
	mob.ZoneId = z.Id
	z.MobileDefinitions[mob.Id] = mob
}

func (z *Zone) AddCommand(c ZoneCommand) {
	z.Commands = append(z.Commands, c)
}

func (z *Zone) String() string {
	return fmt.Sprintf("(Zone %s: '%s')", z.Id, z.Name)
}

func (z *Zone) Reset(mobileRoomMap *MobileRoomMap) []error {
	log.Debug().Str("zone", z.Name).Msg("reset")
	errs := []error{}
	for _, cmd := range z.Commands {
		switch cmd.(type) {
		case CreateMobile:
			var err error
			if err = z.createMobile(mobileRoomMap, cmd.(CreateMobile)); err != nil {
				errs = append(errs, err)
			}
		case CreateObject:
			var err error
			if err = z.createObject(cmd.(CreateObject)); err != nil {
				errs = append(errs, err)
			}
		default:
			errs = append(errs, fmt.Errorf("zone %s unhandled Zone Command Type: %s", z.Id, cmd))
		}

	}

	// Set the last reset time, even if there were errors. Whatever
	// was wrong is probably still wrong the next time, no reason to
	// start spinning errors here.
	z.LastReset = time.Now()

	return errs
}

// Create a mobile.
func (z *Zone) createMobile(mobileRoomMap *MobileRoomMap, cmd CreateMobile) error {
	// TODO determine how many of the definition are in the zone
	defn := z.MobileDefinitions[cmd.MobileDefinitionId]
	if defn == nil {
		return errors.New(fmt.Sprintf("createMobile: definition id not found: %s", cmd))
	}
	// how many of this mobile definition id are in the zone?
	log.Printf("createMobile: considering %s, there are %d, max %d",
		defn.Id, mobileRoomMap.GetMobileDefinitionCount(defn.Id), cmd.InstanceMax)
	if mobileRoomMap.GetMobileDefinitionCount(defn.Id) < cmd.InstanceMax {
		r := z.Rooms[cmd.RoomId]
		if r == nil {
			return errors.New(fmt.Sprintf("createMobile: room not found: %s", cmd))
		}
		log.Printf("Creating mobile: %s", defn.Id)
		mobileRoomMap.Add(mobile.NewInstance(defn), r)
	}
	return nil
}

// Create an object and put it in a room. If there was an error, return the error.
func (z *Zone) createObject(cmd CreateObject) error {
	defn := z.ObjectDefinitions[cmd.ObjectDefinitionId]
	if defn == nil {
		return errors.New(fmt.Sprintf("createObject: definition id not found: %s", cmd))
	}
	r := z.Rooms[cmd.RoomId]
	if r == nil {
		return errors.New(fmt.Sprintf("createObject: room not found: %s", cmd))
	}

	// Note that the amount only deals with how many are in the room (on the floor
	// so to speak) - not lying around attached to mobs, other players, and so on.
	count := 0
	// TODO not the most efficient way of figuring this out ..
	for _, inst := range r.GetAllInventory() {
		if inst.Definition.Identifier() == cmd.ObjectDefinitionId {
			// found one
			count++
		}
	}
	if count >= cmd.InstanceMax {
		log.Printf("Room %s has %d %s, max is %d (not creating more objects)",
			r.Id,
			count, cmd.ObjectDefinitionId,
			cmd.InstanceMax)
		return nil
	}
	return r.AddInventory(object.NewInstance(defn))
}
