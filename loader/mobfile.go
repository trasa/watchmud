package loader

// mob file might not exist, they are optional
type mobEntry struct {
	Id                  string         `json:"id"`
	Name                string         `json:"name"`
	Aliases             []string       `json:"aliases"`
	ShortDescription    string         `json:"short_description"`
	DescriptionInRoom   string         `json:"description_in_room"`
	WanderingDefinition WanderingEntry `json:"wandering_definition"`
	Flags               []string       `json:"flags"`
	MaxHealth           int64          `json:"max_health"`
	AC                  int            `json:"ac"`
}

type WanderingEntry struct {
	CanWander             bool     `json:"can_wander"`
	CheckFrequencySeconds int      `json:"check_frequency_seconds"`
	CheckPercentage       int      `json:"check_percentage"`
	WanderStyle           int      `json:"wander_style"`
	Path                  []string `json:"path"`
}
