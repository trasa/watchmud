package player

import (
	"errors"
	"uuid"

	"github.com/rs/zerolog/log"
	"github.com/trasa/watchmud-message/slot"
	"github.com/trasa/watchmud/object"
	"github.com/trasa/watchmud/rules"
)

type Record struct {
	Id                     uuid.UUID
	Name                   string
	CurHealth, MaxHealth   int64
	LineageId, ClassId     string // string ids, not the int32s (see Phase 6)
	LastZoneId, LastRoomId string
	Abilities              rules.Abilities
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
	lineage, found := cat.Lineages[rec.LineageId]
	if !found {
		log.Error().Str("playerName", rec.Name).Msgf("playerName %s lineage %s not found in catalog", rec.Name, rec.LineageId)
		return nil, errors.New("bad lineage")
	}
	class, found := cat.Classes[rec.ClassId]
	if !found {
		log.Error().Str("playerName", rec.Name).Msgf("playerName %s class %s not found in catalog", rec.Name, rec.ClassId)
		return nil, errors.New("bad class")
	}
	p := &Player{
		Id:        rec.Id,
		Name:      rec.Name,
		out:       out,
		Lineage:   lineage,
		Class:     class,
		inventory: NewInventory(),
		slots:     NewSlots(),
		curHealth: rec.CurHealth,
		maxHealth: rec.MaxHealth,
		abilities: rec.Abilities,
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
			log.Warn().Str("player", rec.Name).Msgf("instance %s not found for slot %s", s.InstanceId, s.Location)
		}
		p.Slots().Set(s.Location, item)
	}
	return p, nil
}

func (p *Player) Record() *Record {
	return &Record{
		Id:         p.Id,
		Name:       p.Name,
		CurHealth:  p.curHealth,
		MaxHealth:  p.maxHealth,
		LineageId:  p.Lineage.Id,
		ClassId:    p.Class.Id,
		Abilities:  p.abilities,
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
