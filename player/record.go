package player

import (
	"errors"
	"uuid"

	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/rules"
	"github.com/trasa/watchmud/slot"
)

type Record struct {
	Id                     uuid.UUID
	Name                   string
	CurHealth, MaxHealth   int64
	LineageId              string // cosmetic; there is no ClassId beside it any more
	LastZoneId, LastRoomId string
	Slots                  []SlotRecord
	Inventory              []InventoryRecord
}

type SlotRecord struct {
	Location   slot.Location
	InstanceId uuid.UUID
}

type InventoryRecord struct {
	InstanceId   uuid.UUID
	ZoneId       string
	DefinitionId string
}

type DefinitionSource interface {
	ObjectDefinition(zoneId, definitionId string) (*object.Definition, bool)
}

func FromRecord(rec *Record, out Sender, cat *rules.Catalog, defs DefinitionSource) (*Player, error) {
	// Same reasoning as the missing definitions below: a lineage that content
	// no longer defines used to refuse the login outright. It is cosmetic now
	// -- it grants nothing and nothing depends on it -- so a retired lineage
	// is worth a log line and a fallback, never a locked-out player.
	lineage, found := cat.Lineages[rec.LineageId]
	if !found {
		log.Warn().Str("player", rec.Name).Msgf("lineage %q not in catalog, falling back to the default", rec.LineageId)
		lineage = cat.DefaultLineage()
		if lineage == nil {
			return nil, errors.New("no lineages defined in the catalog")
		}
	}
	p := &Player{
		id:        rec.Id,
		name:      rec.Name,
		out:       out,
		Lineage:   lineage,
		inventory: NewInventory(),
		slots:     NewSlots(),
		curHealth: rec.CurHealth,
		maxHealth: rec.MaxHealth,
		location:  NewLocation(rec.LastZoneId, rec.LastRoomId),
	}

	for _, ir := range rec.Inventory {
		//Missing definitions. A saved InventoryRecord can reference a zone or object id that content no longer defines — you edit content/, and last week's save now points at nothing. Erroring means an unlucky
		//content edit locks a player out of the game permanently. Log and skip the item so the login succeeds; that's the behavior you want at 2am.
		d, found := defs.ObjectDefinition(ir.ZoneId, ir.DefinitionId)
		if !found {
			log.Warn().Str("player", rec.Name).Msgf("definition not found for %s / %s, dropping it", ir.ZoneId, ir.DefinitionId)
			continue
		}
		i := object.NewInstance(ir.InstanceId, d)
		p.inventory.Add(i)
	}

	for _, s := range rec.Slots {
		item, exists := p.inventory.ByInstanceId(s.InstanceId)
		if !exists {
			// the item didn't survive a content edit; leave the slot empty
			// rather than filling it with a nil that reads as occupied.
			log.Warn().Str("player", rec.Name).Msgf("instance %s not found for slot %s", s.InstanceId, s.Location)
			continue
		}
		p.Slots().Set(s.Location, item)
	}
	return p, nil
}

func (p *Player) Record() *Record {
	return &Record{
		Id:         p.Id(),
		Name:       p.Name(),
		CurHealth:  p.curHealth,
		MaxHealth:  p.maxHealth,
		LineageId:  p.Lineage.Id,
		LastZoneId: p.location.ZoneId,
		LastRoomId: p.location.RoomId,
		Slots:      p.slots.Record(),
		Inventory:  p.inventory.Record(),
	}
}

func (i *Inventory) Record() []InventoryRecord {
	var records []InventoryRecord
	for _, item := range i.GetAll() {
		r := InventoryRecord{
			InstanceId:   item.Id,
			ZoneId:       item.Definition.ObjectId.ZoneId,
			DefinitionId: item.Definition.ObjectId.DefinitionId,
		}
		records = append(records, r)
	}
	return records
}

func (s *Slots) Record() []SlotRecord {
	var records []SlotRecord
	for loc, item := range s.slotMap {
		if item == nil {
			continue
		}
		records = append(records, SlotRecord{
			Location:   loc,
			InstanceId: item.Id,
		})
	}
	return records
}
