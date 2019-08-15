package db

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
)

type ClassData struct {
	ClassId           int32                  `db:"class_id"`
	ClassName         string                 `db:"class_name"`
	AbilityPreference *AbilityPreferenceList `db:"ability_preference" json:"a"`
}

type AbilityPreferenceList struct {
	Preferences []string
}

func (a *AbilityPreferenceList) Value() (driver.Value, error) {
	return json.Marshal(a)
}
func (a *AbilityPreferenceList) Scan(src interface{}) error {
	v := reflect.ValueOf(src)
	if !v.IsValid() || v.CanAddr() && v.IsNil() {
		return nil
	}
	switch ts := src.(type) {
	case []byte:

		return json.Unmarshal(ts, &a)

	case string:
		return json.Unmarshal([]byte(ts), &a)

	default:
		return fmt.Errorf("could not decode type %T -> %T", src, a)
	}
}

func GetClassData() (result []ClassData, err error) {
	err = watchdb.Select(&result, "select c.class_id, c.class_name, c.ability_preference from classes c order by c.class_id")
	return
}

func GetClassDataJson() (classjson []byte, err error) {
	classes, err := GetClassData()
	if err != nil {
		return
	}
	classjson, err = json.Marshal(classes)
	return
}
