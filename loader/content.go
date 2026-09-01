package loader

import (
	"fmt"
	"io/fs"
	"maps"
	"path"
	"slices"
	"time"

	"github.com/trasa/watchmud-message/direction"
	"github.com/trasa/watchmud-message/slot"
	"github.com/trasa/watchmud/behavior"
	"github.com/trasa/watchmud/mobile"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/rules"
	"github.com/trasa/watchmud/spaces"
	"github.com/trasa/watchmud/zonereset"
)

// Content is the static game content read from the content directory and files: everything
// defined on disk and read-only once the server is running.
type Content struct {
	Zones    map[string]*spaces.Zone
	Settings *Settings
	Catalog  *rules.Catalog
}

func NewContent(settings *Settings, catalog *rules.Catalog) *Content {
	return &Content{
		Zones:    make(map[string]*spaces.Zone),
		Settings: settings,
		Catalog:  catalog,
	}
}

func LoadContent(fsys fs.FS) (*Content, error) {
	worldFS, err := fs.Sub(fsys, "world")
	if err != nil {
		return nil, fmt.Errorf("Sub(world): %w", err)
	}

	rulesFS, err := fs.Sub(fsys, "rules")
	if err != nil {
		return nil, fmt.Errorf("Sub(rules): %w", err)
	}
	settings, err := LoadSettings(worldFS)
	if err != nil {
		return nil, err
	}

	cat, err := LoadRulesCatalog(rulesFS)
	if err != nil {
		return nil, err
	}
	c := NewContent(settings, cat)

	if err := c.loadZoneManifest(worldFS); err != nil {
		return nil, err
	}
	if err := c.loadRooms(worldFS); err != nil {
		return nil, err
	}
	if err := c.loadObjectDefinitions(worldFS); err != nil {
		return nil, err
	}
	if err := c.loadMobileDefinitions(worldFS); err != nil {
		return nil, err
	}
	if err := c.loadZoneInstructions(fsys); err != nil {
		return nil, err
	}
	return c, nil
}

func loadInto[T any](dst *T, fsys fs.FS, name string) error {
	v, err := readJSONFile[T](fsys, name)
	if err != nil {
		return err
	}
	*dst = v
	return nil
}

func (c *Content) Room(zoneId, roomId string) (*spaces.Room, error) {
	z, ok := c.Zones[zoneId]
	if !ok {
		return nil, fmt.Errorf("zone %q not found", zoneId)
	}
	r, ok := z.Rooms[roomId]
	if !ok {
		return nil, fmt.Errorf("room %q not found in zone %q", roomId, zoneId)
	}
	return r, nil
}

// Retrieve the zone manifest; prepare the zone objects to be populated by
// rooms, objects, mobiles (but don't process the zone commands yet).
func (c *Content) loadZoneManifest(fsys fs.FS) error {
	manifests, err := readJSONFile[[]zoneManifestEntry](fsys, "zone_manifest.json")
	if err != nil {
		return err
	}
	for _, m := range manifests {
		if !m.Enabled {
			continue
		}
		c.addZone(spaces.NewZone(
			m.Id,
			m.Name,
			zonereset.Mode(m.ResetMode),
			time.Duration(m.LifetimeMinutes)*time.Minute,
		))
	}
	return nil
}

// Read all the room files from all the zones
func (c *Content) loadRooms(fsys fs.FS) error {

	// every room object has to exist before any exit can be connected, so
	// the parsed entries are held until the second pass.
	roomsByZone := make(map[string][]roomFileEntry, len(c.Zones))
	for _, zonename := range c.zoneNames() {
		entries, err := readJSONFile[[]roomFileEntry](fsys, path.Join(zonename, "rooms.json"))
		if err != nil {
			return err
		}
		zone := c.Zones[zonename]

		for i := range entries {
			entry := &entries[i]

			r := spaces.NewRoom(zone, entry.Id, entry.Name, entry.Description)
			r.SetFlags(entry.Flags)
			zone.AddRoom(r)

			// an exit that doesn't name a zone is assumed to stay in this one.
			// Index rather than range-by-value so the write actually sticks.
			for j := range entry.Exits {
				if entry.Exits[j].DestinationZoneId == "" {
					entry.Exits[j].DestinationZoneId = zonename
				}
			}
		}
		roomsByZone[zonename] = entries
	}

	// All rooms exist now, so exits can be resolved
	for _, zonename := range slices.Sorted(maps.Keys(roomsByZone)) {
		for _, entry := range roomsByZone[zonename] {
			for _, exit := range entry.Exits {
				if err := c.connectRooms(
					zonename,
					entry.Id,
					exit.Direction,
					exit.DestinationZoneId,
					exit.DestinationRoomId,
				); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (c *Content) connectRooms(sourceZoneId string, sourceRoomId string, dir direction.Direction, destZoneId string, destRoomId string) error {
	sourceZone := c.Zones[sourceZoneId]
	if sourceZone == nil {
		return fmt.Errorf("connect rooms: source zone %q not found", sourceZoneId)
	}
	destZone := c.Zones[destZoneId]
	if destZone == nil {
		return fmt.Errorf("connect rooms: destination zone %q not found (from %s/%s going %s)",
			destZoneId, sourceZoneId, sourceRoomId, dir)
	}
	sourceRoom := sourceZone.Rooms[sourceRoomId]
	if sourceRoom == nil {
		return fmt.Errorf("connect rooms: source room %q not found in zone %q", sourceRoomId, sourceZoneId)
	}
	destRoom := destZone.Rooms[destRoomId]
	if destRoom == nil {
		return fmt.Errorf("connect rooms: destination room %q not found in zone %q (from %s/%s going %s)",
			destRoomId, destZoneId, sourceZoneId, sourceRoomId, dir)
	}
	sourceRoom.Set(dir, destRoom)
	return nil
}

func (c *Content) addZone(z *spaces.Zone) {
	c.Zones[z.Id] = z
	// TODO should this go away?
}

// zoneNames returns the loaded zone names in a stable order, so that load
// order and error messages are reproducible run to run.
func (c *Content) zoneNames() (result []string) {
	return slices.Sorted(maps.Keys(c.Zones))
}

func (c *Content) loadObjectDefinitions(fsys fs.FS) error {
	// for each zone: create all object definitions
	for _, zonename := range c.zoneNames() {
		objEntries, err := readOptionalJSONFile[[]objectEntry](fsys, path.Join(zonename, "objects.json"))
		if err != nil {
			return err
		}
		for _, obj := range objEntries {
			cat, err := object.StringToCategory(obj.Category)
			if err != nil {
				return fmt.Errorf("object %s/%s: bad category: %w", zonename, obj.Id, err)
			}

			wearLoc, err := slot.StringToLocation(obj.WearLocation)
			if err != nil {
				return fmt.Errorf("object %s/%s: bad wear location: %w", zonename, obj.Id, err)
			}

			d := object.NewDefinition(
				obj.Id,
				obj.Name,
				zonename,
				cat,
				obj.Aliases,
				obj.ShortDescription,
				obj.DescriptionOnGround,
				wearLoc,
			)

			for _, bstr := range obj.Behaviors {
				b, err := behavior.StringToBehavior(bstr)
				if err != nil {
					return fmt.Errorf("object %s/%s: bad behavior %q: %w", zonename, obj.Id, bstr, err)
				}
				d.Behaviors.Add(b)
			}

			c.Zones[zonename].AddObjectDefinition(d)
		}
	}
	return nil
}

func (c *Content) loadMobileDefinitions(fsys fs.FS) error {
	for _, zonename := range c.zoneNames() {
		mobEntries, err := readOptionalJSONFile[[]mobEntry](fsys, path.Join(zonename, "mobs.json"))
		if err != nil {
			return err
		}
		for _, mob := range mobEntries {
			defn := mobile.NewDefinition(
				mob.Id,
				mob.Name,
				zonename,
				mob.Aliases,
				mob.ShortDescription,
				mob.DescriptionInRoom,
				mob.MaxHealth,
				mobile.WanderingDefinition{
					CanWander:       mob.WanderingDefinition.CanWander,
					CheckFrequency:  time.Second * time.Duration(mob.WanderingDefinition.CheckFrequencySeconds),
					CheckPercentage: float32(mob.WanderingDefinition.CheckPercentage) / 100.0,
					Style:           mobile.WanderingStyle(mob.WanderingDefinition.WanderStyle),
					Path:            mob.WanderingDefinition.Path,
				},
				mob.AC,
			)
			defn.SetFlags(mob.Flags)
			c.Zones[zonename].AddMobileDefinition(defn)
		}
	}
	return nil
}

func (c *Content) loadZoneInstructions(fsys fs.FS) error {
	for _, zonename := range c.zoneNames() {
		// instruction files are optional
		insts, err := readOptionalJSONFile[[]instructionFileEntry](fsys, path.Join(zonename, "instructions.json"))
		if err != nil {
			return err
		}
		zone := c.Zones[zonename]
		for _, entry := range insts {
			switch entry.Type {
			case "CreateObject":
				zone.AddCommand(spaces.CreateObject{
					ObjectDefinitionId: entry.ObjectId,
					RoomId:             entry.RoomId,
					InstanceMax:        entry.InstanceMax,
				})
			case "CreateMobile":
				zone.AddCommand(spaces.CreateMobile{
					MobileDefinitionId: entry.MobileId,
					RoomId:             entry.RoomId,
					InstanceMax:        entry.InstanceMax,
				})
			default:
				return fmt.Errorf("zone %s: Unhandled Instruction type: %q", zonename, entry.Type)
			}
		}
	}
	return nil
}
