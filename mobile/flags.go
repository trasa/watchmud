package mobile

// Flag describes various aspects of mob definitions, things that a mob "is".
type Flag string

const (
	Aggressive      Flag = "aggressive"      // aggro, attack players on sight!
	PlayerCantFight Flag = "playerCantFight" // players not allowed to attack this mob
)

func ConvertFlags(flags []string) []Flag {
	var converted []Flag
	for _, flag := range flags {
		converted = append(converted, Flag(flag))
	}
	return converted
}
