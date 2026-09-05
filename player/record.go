package player

import (
	"github.com/google/uuid"
	"github.com/trasa/watchmud/rules"
)

type Record struct {
	Id                     int64
	Name                   string
	CurHealth, MaxHealth   int64
	LineageId, ClassId     string // string ids, not the int32s (see Phase 6)
	LastZoneId, LastRoomId string
	Abilities              rules.Abilities
	Slots                  []SlotRecord
	Inventory              []InventoryRecord
}

type SlotRecord struct {
	Location     int32
	InstanceId   uuid.UUID
	ZoneId       string
	DefinitionId string
}

type InventoryRecord struct {
	InventoryId  int32
	PlayerId     int64
	InstanceId   uuid.UUID
	ZoneId       string
	DefinitionId string
}
