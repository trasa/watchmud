package loader

// object file is optional
type objectEntry struct {
	Id                  string   `json:"id"`
	Name                string   `json:"name"`
	Category            string   `json:"category"`
	Aliases             []string `json:"aliases"`
	ShortDescription    string   `json:"short_description"`
	DescriptionOnGround string   `json:"description_on_ground"`
	WearLocation        string   `json:"wear_location"`
	Behaviors           []string `json:"behaviors"`

	// Roles is what this object contributes to each role while equipped,
	// keyed on rules.Role.Id: {"tank": 3}. Optional, and meaningless on
	// anything without a wear_location.
	Roles map[string]int `json:"roles"`
}
