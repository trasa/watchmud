package db

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type ClassData struct {
	ClassId           int32                  `db:"class_id"`
	ClassName         string                 `db:"class_name"`
	AbilityPreference *AbilityPreferenceList `db:"ability_preference" json:"a"`
}

type AbilityPreferenceList struct {
	Preferences []string `json:"a"`
}

func (a AbilityPreferenceList) Value() (driver.Value, error) {
	return json.Marshal(a)
}
func (a *AbilityPreferenceList) Scan(src interface{}) error {
	bs, ok := src.([]byte)
	if !ok {
		return errors.New("not a []byte")
	}
	return json.Unmarshal(bs, a)
}

func GetClassData() (result []ClassData, err error) {
	err = watchdb.Select(&result, "select c.class_id, c.class_name, c.ability_preference from classes c order by c.class_id")
	return
}

func GetSingleClassData(classId int32) (result ClassData, err error) {
	err = watchdb.Get(&result, "select c.class_id, c.class_name, c.ability_preference from classes c where c.class_id = $1", classId)
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
