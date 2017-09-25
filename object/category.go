package object

import "strings"

type Category int

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

func (cs CategorySet) ToList() (result []Category) {
	for k, v := range cs {
		if v {
			result = append(result, k)
		}
	}
	return
}

func (cs CategorySet) Add(c Category) {
	cs[c] = true
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
