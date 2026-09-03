package rules

import "fmt"

type Class struct {
	Id                string   `json:"id"`
	Name              string   `json:"name"`
	AbilityPreference []string `json:"ability_preference" json:"a"`
}

func indexClasses(all []*Class) (map[string]*Class, error) {
	classes := make(map[string]*Class, len(all))
	for _, c := range all {
		if c.Id == "" {
			return nil, fmt.Errorf("class %q: missing id", c.Name)
		}
		if prev, dup := classes[c.Id]; dup {
			return nil, fmt.Errorf("duplicate class id %q (%s and %s)", c.Id, prev.Name, c.Name)
		}
		// validate that the ability preferences are valid attributes, if present
		for _, attr := range c.AbilityPreference {
			if attr != "str" && attr != "dex" && attr != "con" && attr != "int" && attr != "wis" && attr != "cha" {
				return nil, fmt.Errorf("class %q: invalid ability preference %q", c.Name, attr)
			}
		}
		classes[c.Id] = c
	}
	return classes, nil
}
