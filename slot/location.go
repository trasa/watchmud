package slot

type Location int32

//go:generate stringer -type=Location
const (
	None  Location = iota
	Wield          // item can be wielded (i.e. weapon)
	Hold           // item can be held ('hold' command)
)
