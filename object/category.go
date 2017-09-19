package object

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
