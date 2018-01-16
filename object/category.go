package object

import "strings"

type Category int32

//go:generate stringer -type=Category
const (
	WEAPON Category = iota
	WAND
	STAFF
	TREASURE
	ARMOR
	FOOD
	OTHER
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

func IntsToCategories(rawcats []int32) (result []Category) {
	for _, i := range rawcats {
		result = append(result, Category(i))
	}
	return
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
