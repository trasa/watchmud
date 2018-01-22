package slot

type Location int32

//go:generate stringer -type=Location
const (
	None Location = iota
	PrimaryHand
	SecondaryHand
)
