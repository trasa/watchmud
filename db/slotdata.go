package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/trasa/watchmud/object"
	"reflect"
)

type SlotDataList struct {
	Slots []SlotData `json:"Slots"`
}

type SlotData struct {
	Location     int32     `json:"Location"`
	InstanceId   uuid.UUID `json:"InstanceId"`
	ZoneId       string    `json:"ZoneId"`
	DefinitionId string    `json:"DefinitionId"`
}

func NewSlotDataList(slots *object.Slots) (result *SlotDataList) {
	result = &SlotDataList{}
	for location, inst := range slots.GetAll() {
		result.Slots = append(result.Slots, SlotData{
			Location:     int32(location),
			InstanceId:   inst.InstanceId,
			ZoneId:       inst.Definition.ZoneId,
			DefinitionId: inst.Definition.Identifier(),
		})
	}
	return result
}

// driver.Valuer
func (sd *SlotDataList) Value() (driver.Value, error) {
	return json.Marshal(sd)
}

// sql.Scanner
func (sd *SlotDataList) Scan(src interface{}) error {
	v := reflect.ValueOf(src)
	if !v.IsValid() || v.CanAddr() && v.IsNil() {
		return nil
	}
	switch ts := src.(type) {
	case []byte:

		return json.Unmarshal(ts, &sd)

	case string:
		return json.Unmarshal([]byte(ts), &sd)

	default:
		return fmt.Errorf("could not decode type %T -> %T", src, sd)
	}
}

// driver.Valuer
func (sd *SlotData) Value() (driver.Value, error) {
	return json.Marshal(sd)
}

// sql.Scanner
func (sd *SlotData) Scan(src interface{}) error {
	v := reflect.ValueOf(src)
	if !v.IsValid() || v.CanAddr() && v.IsNil() {
		return nil
	}
	switch ts := src.(type) {
	case []byte:

		return json.Unmarshal(ts, &sd)

	case string:
		return json.Unmarshal([]byte(ts), &sd)

	default:
		return fmt.Errorf("could not decode type %T -> %T", src, sd)
	}
}
