package rules

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

type ObjectCategory string

const (
	ObjectCategoryNone     ObjectCategory = ""
	ObjectCategoryWeapon   ObjectCategory = "weapon"
	ObjectCategoryWand     ObjectCategory = "wand"
	ObjectCategoryStaff    ObjectCategory = "staff"
	ObjectCategoryTreasure ObjectCategory = "treasure"
	ObjectCategoryArmor    ObjectCategory = "armor"
	ObjectCategoryFood     ObjectCategory = "food"
	ObjectCategoryOther    ObjectCategory = "other"
	ObjectCategoryCorpse   ObjectCategory = "corpse"
)

var ErrUnknownObjectCategory = errors.New("unknown object category")

var objectCategories = []ObjectCategory{
	ObjectCategoryNone,
	ObjectCategoryWeapon,
	ObjectCategoryWand,
	ObjectCategoryStaff,
	ObjectCategoryTreasure,
	ObjectCategoryArmor,
	ObjectCategoryFood,
	ObjectCategoryOther,
	ObjectCategoryCorpse,
}

func (c *ObjectCategory) UnmarshalText(b []byte) error {
	s, err := ParseObjectCategory(string(b))
	if err != nil {
		return err
	}
	*c = s
	return nil
}

func ParseObjectCategory(s string) (ObjectCategory, error) {
	c := ObjectCategory(strings.ToLower(s))
	if !slices.Contains(objectCategories, c) {
		return ObjectCategoryNone, fmt.Errorf("%w: %q", ErrUnknownObjectCategory, s)
	}
	return c, nil
}
