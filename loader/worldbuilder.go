package loader

import (
	"errors"
	"fmt"
	"github.com/trasa/watchmud/direction"
	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/spaces"
	"github.com/trasa/watchmud/zonereset"
	"path/filepath"
	"time"
)

type WorldBuilder struct {
	Zones         map[string]*spaces.Zone
	Settings      *SettingsFile
	worldFilesDir string
}

// Build the world, using the files found under the worldFilesDirectory
func BuildWorld(worldFilesDirectory string) (*WorldBuilder, error) {
	worldBuilder := &WorldBuilder{
		Zones:         make(map[string]*spaces.Zone),
		worldFilesDir: worldFilesDirectory,
	}

	if err := worldBuilder.loadSettings(); err != nil {
		return nil, err
	}

	if err := worldBuilder.loadZoneManifest(); err != nil {
		return nil, err
	}

	if err := worldBuilder.loadRooms(); err != nil {
		return nil, err
	}

	if err := worldBuilder.loadObjectDefinitions(); err != nil {
		return nil, err
	}

	if err := worldBuilder.loadMobileDefinitions(); err != nil {
		return nil, err
	}

	if err := worldBuilder.loadZoneInstructions(); err != nil {
		return nil, err
	}

	return worldBuilder, nil
}

func (wb *WorldBuilder) loadSettings() error {
	settings, err := readSettingsFile(filepath.Join(wb.worldFilesDir, "settings.json"))
	if err != nil {
		return err
	}
	wb.Settings = settings
	return nil
}

// Retrieve the zone manifest; prepare the zone objects to be
// populated by rooms, objects, mobiles (but don't process the
// zone commands yet)
func (wb *WorldBuilder) loadZoneManifest() error {

	zonemanifests, err := readZoneManifest(filepath.Join(wb.worldFilesDir, "zone_manifest.json"))
	if err != nil {
		return err
	}
	for _, manifest := range zonemanifests {
		if manifest.Enabled {
			z := spaces.NewZone(manifest.Id, manifest.Name, zonereset.Mode(manifest.ResetMode), time.Minute*time.Duration(manifest.LifetimeMinutes))
			wb.addZone(z)
		}
	}
	return nil
}

// Read all the room files from all the zones
func (wb *WorldBuilder) loadRooms() error {

	// have to create all the room objects before we can
	// create all the connections - so must store the
	// information for later
	roomMap := make(map[string][]roomFileEntry)

	for _, zonename := range wb.zoneNames() {
		fileEntries, err := readRoomFile(filepath.Join(wb.worldFilesDir, zonename, "rooms.json"))
		if err != nil {
			return err
		}
		roomMap[zonename] = fileEntries
		for _, roomEntry := range fileEntries {
			r := spaces.NewRoom(wb.Zones[zonename], roomEntry.Id, roomEntry.Name, roomEntry.Description)
			wb.Zones[zonename].AddRoom(r)
			// If the direction zone/room doesn't indicate the zone,
			// assume that the direction is to a zoome in the current zone
			for _, exitInfo := range roomEntry.Exits {
				if exitInfo.DestinationZoneId == "" {
					exitInfo.DestinationZoneId = zonename
				}
			}
		}
	}

	// all the rooms are created, now can loop over and create
	// all the connections
	for zonename, roomEntries := range roomMap {
		for _, roomEntry := range roomEntries {
			for _, exitInfo := range roomEntry.Exits {
				if err := wb.connectRooms(zonename,
					roomEntry.Id,
					direction.Direction(exitInfo.Direction),
					exitInfo.DestinationZoneId,
					exitInfo.DestinationRoomId); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (wb *WorldBuilder) connectRooms(sourceZoneId string, sourceRoomId string, dir direction.Direction, destZoneId string, destRoomId string) error {
	sourceZone := wb.Zones[sourceZoneId]
	if sourceZone == nil {
		return errors.New(fmt.Sprintf("connectRooms failed: sourceZoneId '%s' not found", sourceZoneId))
	}
	destZone := wb.Zones[destZoneId]
	if destZone == nil {
		return errors.New(fmt.Sprintf("connectRooms failed: destZoneId '%s' not found", destZoneId))
	}
	sourceRoom := sourceZone.Rooms[sourceRoomId]
	if sourceRoom == nil {
		return errors.New(fmt.Sprintf("connectRooms failed: zone %s sourceRoomId '%s' not found", sourceZoneId, sourceRoomId))
	}
	destRoom := destZone.Rooms[destRoomId]
	if destRoom == nil {
		return errors.New(fmt.Sprintf("ConnectRooms failed: zone %s destRoomId '%s' not found", destZoneId, destRoomId))
	}
	sourceRoom.Set(dir, destRoom)
	return nil
}

func (wb *WorldBuilder) addZone(z *spaces.Zone) {
	wb.Zones[z.Id] = z
}

func (wb *WorldBuilder) zoneNames() (result []string) {
	for k := range wb.Zones {
		result = append(result, k)
	}
	return
}

func (wb *WorldBuilder) loadObjectDefinitions() error {

	// for each zone: create all object definitions
	for _, zonename := range wb.zoneNames() {
		objEntries, err := readObjectFile(filepath.Join(wb.worldFilesDir, zonename, "objects.json"))
		if err != nil {
			return err
		}
		for _, obj := range objEntries {
			wb.Zones[zonename].AddObjectDefinition(
				object.NewDefinition(obj.Id,
					obj.Name,
					zonename,
					object.Category(obj.Category),
					obj.Aliases,
					obj.ShortDescription,
					obj.DescriptionOnGround))
		}
	}
	return nil
}

func (wb *WorldBuilder) loadMobileDefinitions() error {

	for _, zonename := range wb.zoneNames() {
		mobEntries, err := readMobFile(filepath.Join(wb.worldFilesDir, zonename, "mobs.json"))
		if err != nil {
			return err
		}
		for _, mob := range mobEntries {
			wb.Zones[zonename].AddMobileDefinition(
				mobile.NewDefinition(mob.Id,
					mob.Name,
					zonename,
					mob.Aliases,
					mob.ShortDescription,
					mob.DescriptionInRoom,
					mobile.WanderingDefinition{
						CanWander:       mob.WanderingDefinition.CanWander,
						CheckFrequency:  time.Second * time.Duration(mob.WanderingDefinition.CheckFrequencySeconds),
						CheckPercentage: float32(mob.WanderingDefinition.CheckPercentage) / 100.0,
						Style:           mobile.WanderingStyle(mob.WanderingDefinition.WanderStyle),
						Path:            mob.WanderingDefinition.Path,
					}))
		}
	}
	return nil
}

func (wb *WorldBuilder) loadZoneInstructions() error {

	for _, zonename := range wb.zoneNames() {
		insts, err := readInstructionFile(filepath.Join(wb.worldFilesDir, zonename, "instructions.json"))
		if err != nil {
			return err
		}
		for _, entry := range insts {
			switch entry.Type {
			case "CreateObject":
				wb.Zones[zonename].AddCommand(spaces.CreateObject{
					ObjectDefinitionId: entry.ObjectId,
					RoomId:             entry.RoomId,
					InstanceMax:        entry.InstanceMax,
				})
			case "CreateMobile":
				wb.Zones[zonename].AddCommand(spaces.CreateMobile{
					MobileDefinitionId: entry.MobileId,
					RoomId:             entry.RoomId,
					InstanceMax:        entry.InstanceMax,
				})
			default:
				return errors.New(fmt.Sprintf("Unhandled Instruction type: %s", entry.Type))
			}
		}
	}
	return nil
}
