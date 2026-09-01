package loader

// object file is optional
type objectFileEntry struct {
	Id                  string   `json:"id"`
	Name                string   `json:"name"`
	Category            string   `json:"category"`
	Aliases             []string `json:"aliases"`
	ShortDescription    string   `json:"short_description"`
	DescriptionOnGround string   `json:"description_on_ground"`
	WearLocation        string   `json:"wear_location"`
	Behaviors           []string `json:"behaviors"`
}
