package rules

type Class struct {
	Id                string   `json:"id"`
	Name              string   `json:"name"`
	AbilityPreference []string `json:"ability_preference" json:"a"`
}
