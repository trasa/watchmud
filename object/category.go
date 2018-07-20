package object

import (
	"fmt"
	"strings"
)

type Category int32

//go:generate stringer -type=Category
const (
	None Category = iota
	Weapon
	Wand
	Staff
	Treasure
	Armor
	Food
	Other
	Corpse
)

type CategorySet map[Category]bool

func (cs CategorySet) ToInt32List() (result []int32) {
	for k, v := range cs {
		if v {
			result = append(result, int32(k))
		}
	}
	return
}

func (cs CategorySet) Add(c Category) {
	cs[c] = true
}

func (cs CategorySet) Contains(c Category) bool {
	return cs[c]
}

func CategoriesToString(cats []Category) string {
	if len(cats) == 0 {
		return ""
	} else {
		var strs []string
		for _, c := range cats {
			strs = append(strs, c.String())
		}
		return strings.Join(strs, ", ")
	}
}

func StringToCategory(categoryName string) (Category, error) {
	if categoryName == "" {
		return None, nil
	}
	categoryName = strings.ToUpper(categoryName)
	stridx := strings.Index(strings.ToUpper(_Category_name), categoryName)
	if stridx < 0 {
		return None, fmt.Errorf("category '%s' not found", categoryName)
	}

	for pos, catidx := range _Category_index {
		if stridx == int(catidx) {
			return Category(pos), nil
		}
	}
	// shouldn't happen?
	return None, fmt.Errorf("could not find index %d for category '%s'", stridx, categoryName)
}
