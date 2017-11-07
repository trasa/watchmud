package zonereset

// When to reset the zone?
type Mode int

const (
	NEVER      Mode = iota // Don't ever reset. Obviously.
	NO_PLAYERS             // Reset possible only when the zone is empty of players.
	ALWAYS                 // Reset at the correct time, no matter what else is happening.
)
