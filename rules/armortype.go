package rules

import (
	"errors"
	"fmt"
	"slices"
)

type ArmorType string

const (
	ArmorTypeNone    ArmorType = ""
	ArmorTypeCloth   ArmorType = "cloth"
	ArmorTypeLeather ArmorType = "leather"
	ArmorTypePlate   ArmorType = "plate"
)

func (t *ArmorType) UnmarshalText(b []byte) error {
	armorType, err := ParseArmorType(string(b))
	if err != nil {
		return err
	}
	*t = armorType
	return nil
}

var ErrUnknownArmorType = errors.New("unknown armor type")

var armorTypes = []ArmorType{
	ArmorTypeCloth,
	ArmorTypeLeather,
	ArmorTypePlate,
}

func ArmorTypes() []ArmorType {
	return slices.Clone(armorTypes)
}

func ParseArmorType(s string) (ArmorType, error) {
	at := ArmorType(s)
	if !slices.Contains(armorTypes, at) {
		return ArmorTypeNone, fmt.Errorf("%w: %q", ErrUnknownArmorType, s)
	}
	return at, nil
}
